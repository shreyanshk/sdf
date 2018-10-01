==================
SDF: Sane DotFiles
==================

SDF allows setting up your dotfiles as simple as:

.. code-block:: console

	$ sdf clone <URL to your repository>
	$ sdf checkout master

What is SDF?
============

SDF is a program to help you version control your $HOME directory.

Under the hood it wraps ``git`` to make version control of your dotfiles simple and straightforward while also giving you the full flexibility of Git's CLI when you need it. It reimplements a few commands to make them more suitable for handling the task of managing dotfiles and then provides some useful extra commands.

Installing SDF
==============

SDF is compiled with `Go`_ compiler. It uses `Git`_ internally and optionally depends on `strace`_ for tracing (see below).

.. _Go: https://github.com/golang/go
.. _strace: https://github.com/strace/strace
.. _Git: https://git-scm.com/

Build + Install
---------------

.. code-block:: console

    $ go build -o sdf main.go   # compile
    $ mv sdf /usr/bin/sdf       # or use directly

Get started with SDF
====================

Start with a clean $HOME
------------------------

Create an empty repository on any Git hosting service. Let's assume your repository URL is ``https://example.com/username/profile.git``.

1. Initialize SDF like this:

.. code-block:: console

    $ sdf init https://example.com/username/profile.git
    Initialized new configuration.

2. Your first commit and push:

.. code-block:: console

    $ sdf add ~/.bashrc
    $ sdf commit -m "Initial commit"
    $ sdf push --set-upstream origin master

Here, ``master`` is the branch that you're pushing to.

3. Rinse & Repeat

.. code-block:: console

    $ sdf add ~/.zshrc
    $ sdf commit -m "Move to Zsh"
    $ sdf push

Restore previous $HOME
----------------------

Assuming your repository URL is ``https://example.com/username/profile.git`` and you need the ``master`` branch.

.. code-block:: console

	$ sdf clone https://example.com/username/profile.git
	$ sdf checkout master

That's really it!

Advanced usage
==============

Tracing
-------

At times, you don't know (or remember) which configuration files in your $HOME are being read by a program but you want to version control them.

**Presenting "sdf trace" to your disposal.** SDF with the help of venerable `strace`_ will help you version control configuration for the awesome music library manager `beets`_ as an example.

.. _beets: https://github.com/beetbox/beets

.. code-block:: console

    $ sdf trace beet
    .config/beets/config.yaml

You now know that ``$HOME/.config/beets/config.yaml`` is the file you need.

Git CLI
-------

Because SDF is a wrapper around Git, you can pass all valid git commands (except clone & init) like:

.. code-block:: console

   $ sdf checkout -b dev   # switch branch
   $ sdf log               # view changelog
   $ sdf diff @~..@        # view diff of last commit

See `Git's documentation`_ for more details.

.. _Git's documentation: https://git-scm.com/doc

Credits
=======

* Shreyansh Khajanchi <shreyansh_k@live.com>
* SneakyCobra <https://news.ycombinator.com/item?id=11070797>
