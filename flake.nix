{
  inputs = {
    nixpkgs.url = "https://channels.nixos.org/nixos-unstable/nixexprs.tar.xz";
    multiverse.url = "github:fzakaria/nixpkgs-multiverse";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
  };

  outputs =
    inputs@{ flake-parts, multiverse, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      perSystem =
        {
          lib,
          pkgs,
          self',
          system,
          ...
        }:
        let
          mv = inputs.multiverse.multiverse.${system};
          pins = inputs.multiverse.lib.readLock {
            inherit system;
            file = ./multiverse.lock;
          };
        in
        {
          formatter = pkgs.nixfmt-tree;

          packages.mvs = multiverse.packages.${system}.mvs;

          packages.golangci-lint-v2 = pkgs.runCommandLocal "golangci-lint-v2" { } ''
            mkdir -p $out/bin
            cp ${lib.getExe pins.golangci-lint} $out/bin/$name
          '';

          devShells.default = pkgs.mkShell {
            CGO_ENABLED = 0;
            GOTOOLCHAIN = "local";

            nativeBuildInputs =
              builtins.attrValues pins
              ++ (with pkgs; [
                git
                git-lfs

                # core
                go-task
                gofumpt
                gopls
                gotools
                self'.packages.golangci-lint-v2

                # ui
                nodejs_24
                yarn-berry_3
              ]);
          };
        };
    };
}
