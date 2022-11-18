{
  description = "Package for webhook-dispatcher";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in {
        packages = {
          default = self.packages.${system}.webhook-dispatcher;
          webhook-dispatcher = pkgs.buildGoModule {
            name = "webhook-dispatcher";
            src = ./.;
            vendorSha256 = "sha256-K/8AfZQq4Y+nnKqAYFa/+Lojr2r+HkT+MmvMHh43C58=";
          };
        };

        apps = {
          default = self.apps.${system}.webhook-dispatcher;
          webhook-dispatcher = {
            type = "app";
            program =
              "${self.packages.${system}.webhook-dispatcher}/bin/webhook-dispatcher";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gotools go-tools gopls ];
        };
      });
}
