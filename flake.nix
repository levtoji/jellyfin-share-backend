{
  description = "Jellyfin Share Backend - Enhanced with direct links, audio/subtitle selection";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          jellyfin-share-backend = pkgs.callPackage ./package.nix {};
        in {
          inherit jellyfin-share-backend;
          default = jellyfin-share-backend;
        }
      );

      nixosModules.default = import ./module.nix;
    };
}
