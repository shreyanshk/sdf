with import <nixpkgs> {};

buildGoPackage rec {
  name = "sdf-unstable-${version}";
  version = "0.1";

  goPackagePath = "github.com/shreyansh_k/sdf";

  src = ./.;

  # TODO: add metadata https://nixos.org/nixpkgs/manual/#sec-standard-meta-attributes
  meta = {
    description = "Sane dotfiles";
  };
}
