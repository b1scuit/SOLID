{
  description = "Solid project workspace";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let 
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = 
        let
            inherit (pkgs) stdenv lib;
        in
        stdenv.mkDerivation {
          buildInputs = [
            pkgs.go
          ];
          name = "solid";
          src = "./.";
          cleanPhase = "rm -rf $sourceRoot/*";
          unpackPhase = ''
            sourceRoot=$PWD
            mkdir $sourceRoot
            cp -r $src/* $sourceRoot/
          '';
        };
        devShells.default =
        let
            inherit (pkgs) stdenv lib;
        in
        pkgs.mkShell {
            name = "go";
            buildInputs = [
              pkgs.cowsay
              pkgs.lolcat
              pkgs.go
            ];

            shellHook = ''
              echo "Go Shell" | cowsay | lolcat
            '';
          };
      }
      );
}
