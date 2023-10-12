{
  description = "Clipboard manager for NixOS";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        version = "3.2.0" + "-" + (self.shortRev or "dirty");
        pkgs = nixpkgs.legacyPackages.${system};
      in rec {
        packages.default = pkgs.buildGoModule {
          pname = "clipboard";
          inherit version;
          src = ./.;
          #vendorSha256 = pkgs.lib.fakeSha256;
          vendorHash = "sha256-u93OLlHQwmqKXKTo7J5aCyVurQG+j+v6ES6d713kqvk=";
        };
      }
    );
}
