with import <nixpkgs> {};

buildGoPackage rec {
  name = "sdf-unstable-${version}";
  version = "2018-09-30";
  rev = "cebc0c9af1718c58d202f1e3702bcd51dcfb5c01";

  goPackagePath = "github.com/shreyansh_k/sdf";

  src = ./.;

  goDeps = ./deps.nix;

  # TODO: add metadata https://nixos.org/nixpkgs/manual/#sec-standard-meta-attributes
  meta = {
    description = "Sane dotfiles";
  };
}
