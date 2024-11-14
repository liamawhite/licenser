{

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    {
        overlay = final: prev: { inherit (self.packages.${final.system}) licenser; };
    }
    // 
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells = {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
            ];
          };
        };
        packages = rec {
          licenser = pkgs.callPackage ./licenser.nix { };
          default = licenser;
        };
      }
    );
}

