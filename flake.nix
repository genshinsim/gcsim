{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      perSystem =
        {
          lib,
          pkgs,
          self',
          ...
        }:
        {
          formatter = pkgs.nixfmt-tree;

          packages.golangci-lint-v2 =
            with pkgs;
            stdenv.mkDerivation {
              name = "golangci-lint-v2";
              buildInputs = [ golangci-lint ];
              phases = [ "installPhase" ];
              installPhase = ''
                mkdir -p $out/bin
                cp ${lib.getExe golangci-lint} $out/bin/$name
              '';
            };

          devShells.default = pkgs.mkShell {
            nativeBuildInputs =
              with pkgs;
              let
                go = go_1_25;
                nodejs = nodejs_24;
                yarn = yarn-berry_3;
              in
              [
                git
                git-lfs

                # core
                go
                go-task
                gofumpt
                golangci-lint
                self'.packages.golangci-lint-v2
                gopls
                gotools

                # ui
                nodejs
                yarn

                # protobuf
                protobuf
                protoc-gen-go
                protoc-gen-go-grpc
              ];
          };
        };
    };
}
