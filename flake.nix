{
  description = "Minimal, customizable, and neofetch-like weather CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};

        pname = "stormy";
        version = "0.3.2";
      in {
        packages.default = pkgs.buildGoModule {
          inherit pname version;

          src = self;

          vendorHash = "sha256-iwgGAJRygi+xS5eorZ8wyR6pMDZvmGFXBbCiFazyaHw=";

          meta = with pkgs.lib; {
            description = "Minimal, customizable, and neofetch-like weather CLI";
            license = licenses.mit;
            mainProgram = pname;
          };
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/${pname}";
        };

        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
            gopls
            pkg-config
          ];
        };
      }
    );
}
