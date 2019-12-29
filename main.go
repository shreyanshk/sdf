package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// changing XDG_CONFIG_HOME after initialization
// will require a complete reinstall of profile
// TODO can we make it honor XDG_CONFIG_HOME without hassle?
var userPath = os.Getenv("HOME")
var sdfPath = userPath + "/.config/sdf"
var baseGit = "git --git-dir=" + sdfPath +
	" --work-tree=" + userPath

// sdf <valid git command>
// Escape the abstractions! Get full access to the underlying repository.
func delegateCmdToGit(cmd []string) {
	fullCmd := append(strings.Fields(baseGit), cmd...)
	runWithOutput(fullCmd...)
}

// sdf clone <url>
// Initialize the configuration from repository.
func initFromVCS(url string) {
	if isInitialized() {
		fmt.Println("SDF is already initialized.")
		if !askForConfirmation("Force remove previous configuration?") {
			return
		}
		check(os.RemoveAll(sdfPath))
	}
	// Git magic below
	check(os.MkdirAll(userPath+"/.config", 0600))
	tempDir := sdfPath + "-tmp"
	runWithOutput(
		"git", "clone", "--separate-git-dir="+
			sdfPath, url, tempDir,
	)
	// ensure git-modules work.
	modules := tempDir + "/.gitmodules"
	if _, err := os.Stat(modules); !os.IsNotExist(err) {
		check(os.Rename(modules, userPath+"/.gitmodules"))
	}
	check(os.RemoveAll(tempDir))
	gitCmd2 := append(
		strings.Fields(baseGit),
		"config", "status.showUntrackedFiles", "no",
	)
	exec.Command(gitCmd2[0], gitCmd2[1:]...).Run()
	// ensure other users can't see our data.
	check(os.Chmod(sdfPath, 0700))
	fmt.Println("Restored SDF configuration, activate it with 'sdf checkout .'")
}

// sdf init <url>
// Initialize a new configuration and set default remote URL.
func initNew(url string) {
	if isInitialized() {
		fmt.Println("SDF is already initialized!")
		if !askForConfirmation("Force remove previous configuration?") {
			return
		}
		check(os.RemoveAll(sdfPath))
	}
	// Git magic below
	runWithOutput(
		"git", "init", "--bare",
		sdfPath,
	)
	// This block sets the remote URL
	gitCmd1 := append(
		strings.Fields(baseGit),
		"remote", "add", "origin", url,
	)
	exec.Command(gitCmd1[0], gitCmd1[1:]...).Run()
	gitCmd2 := append(
		strings.Fields(baseGit),
		"config", "status.showUntrackedFiles", "no",
	)
	exec.Command(gitCmd2[0], gitCmd2[1:]...).Run()
	// ensure other users can't see our data.
	check(os.Chmod(sdfPath, 0700))
	fmt.Println("Initialized new configuration.")
}

// sdf trace <command>
// Launch the given program under strace and then filters
// output to display the files that are opened by it.
func traceCmd(inCmd []string) {
	// test if strace is present
	if _, err := exec.LookPath("strace"); err != nil {
		fmt.Println("Strace not found. Check your $PATH or install it.")
		return
	}
	// test if given binary exist
	if _, err := exec.LookPath(inCmd[0]); err != nil {
		fmt.Println("Binary not executable or doesn't exist. Cannot continue.")
		return
	}
	straceArgs := strings.Fields("-f -e trace=openat")
	fullArgs := append(straceArgs, inCmd...)
	straceCmd := exec.Command("strace")
	straceCmd.Args = append(straceCmd.Args, fullArgs...)
	straceOut, err := straceCmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewReader(straceOut)
	straceCmd.Start()
	uplen := len(userPath) // needed for cleaning output
	for {
		line, err := scanner.ReadString('\n')
		if err == io.EOF {
			break
		}
		temp := strings.Split(line, "\"")
		if len(temp) > 2 { // make sure line has a valid path
			fullpath := temp[1]                        // extract the path
			if strings.HasPrefix(fullpath, userPath) { // show stuff from $HOME
				if node, err := os.Stat(fullpath); err == nil && !node.IsDir() {
					fmt.Println(fullpath[uplen+1:]) // remove $HOME prefix
				} else if err != nil && !os.IsNotExist(err) { // let user handle other errors
					panic(err)
				}
			}
		}
	}
	straceCmd.Wait() // reap process entry from process table
}

// Bunch of helper functions.

// test if SDF is initialized.
func isInitialized() bool {
	if _, err := os.Stat(sdfPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func runWithOutput(cmdStr ...string) {
	cmd := exec.Command(
		cmdStr[0], cmdStr[1:]...,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const helpstr = `Usage: sdf <command> [<args>]

SDF: Sane DotFiles
Manage your dotfiles with ease.

SDF is a wrapper around Git and helps version control dotfiles.
It reimplements a few commands and provides some more.

SDF specific commands are:

   clone <url>   Clone user profile configuration from given URL
   init  <url>   Create an empty profile with given URL as upstream
   trace <arg>   List files opened by given command during runtime

Because SDF is just a wrapper around Git, you can pass all valid
git commands like so:

   $ sdf add ~/.bashrc
   $ sdf commit -m "Initial commit"

See Git's documentation with 'man git' for more details.
`

func main() {
	// handle no arguments case
	if len(os.Args) == 1 {
		fmt.Printf(helpstr)
		return
	}
	switch os.Args[1] {
	case "clone":
		if len(os.Args) >= 4 {
			fmt.Println("Too many parameters.")
			return
		} else if len(os.Args) != 3 {
			fmt.Println("URL required.")
			return
		}
		initFromVCS(os.Args[2])
		return
	case "init":
		if len(os.Args) >= 4 {
			fmt.Println("Too many parameters.")
			return
		} else if len(os.Args) != 3 {
			fmt.Println("URL required.")
			return
		}
		initNew(os.Args[2])
		return
	case "trace":
		if len(os.Args) < 3 {
			fmt.Println("Please provide command.")
			return
		}
		traceCmd(os.Args[2:])
		return
	default:
		delegateCmdToGit(os.Args[1:])
		return
	}
}
