{
  description = "URL Shortener Nix flake with systemd template module";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  };

  outputs = { self, nixpkgs }:
    let
      lib = nixpkgs.lib;
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: lib.genAttrs systems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in {
      nixosModules.url-shortener = { config, lib, pkgs, ... }:
        let
          cfg = config.services.url-shortener;
        in {
          options.services.url-shortener = {
            enable = lib.mkEnableOption "URL Shortener service";

            binaryPath = lib.mkOption {
              type = lib.types.path;
              description = "Path to the url-shortener Go binary.";
            };

            frontendDir = lib.mkOption {
              type = lib.types.path;
              default = "/var/lib/url-shortener/frontend";
              description = "Path to the static frontend assets directory.";
            };

            defaultPort = lib.mkOption {
              type = lib.types.int;
              default = 8080;
              description = "Default port for the template instance.";
            };

            environmentFileDir = lib.mkOption {
              type = lib.types.str;
              default = "/etc/url-shortener";
              description = "Directory containing environment files for each instance.";
            };

            user = lib.mkOption {
              type = lib.types.str;
              default = "url-shortener";
              description = "User that runs the service.";
            };

            group = lib.mkOption {
              type = lib.types.str;
              default = "url-shortener";
              description = "Group that runs the service.";
            };
          };

          config = lib.mkIf cfg.enable {
            users.users.${cfg.user} = {
              isSystemUser = true;
              group = cfg.group;
            };

            users.groups.${cfg.group} = { };

            systemd.services."url-shortener@" = {
              description = "URL Shortener instance %i";
              after = [ "network.target" ];
              wantedBy = [ "multi-user.target" ];

              serviceConfig = {
                Type = "simple";
                User = cfg.user;
                Group = cfg.group;
                EnvironmentFile = "${cfg.environmentFileDir}/%i.env";
                ExecStart = "${cfg.binaryPath}";
                Restart = "on-failure";
              };

              environment = {
                FRONTEND_DIR = cfg.frontendDir;
                SERVER_PORT = "%i";
              };
            };
          };
        };

      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.nodejs_20
            pkgs.sqlite
            pkgs.pkg-config
          ] ++ lib.optionals pkgs.stdenv.isDarwin [
            pkgs.darwin.libresolv
            pkgs.clang
          ];
          shellHook = ''
            ${lib.optionalString pkgs.stdenv.isDarwin "export CGO_ENABLED=1"}
            ${lib.optionalString pkgs.stdenv.isDarwin "export NIX_LDFLAGS=\\\"-L${pkgs.darwin.libresolv}/lib\\\""}
            ${lib.optionalString pkgs.stdenv.isDarwin "export CGO_LDFLAGS=\\\"-L${pkgs.darwin.libresolv}/lib\\\""}
            echo "URL-Shortener dev shell loaded"
            echo "Optional env vars: CAP_SITEVERIFY_URL, CAP_SECRET, CAP_API_ENDPOINT, SHORTEN_PASSWORD"
          '';
        };
      });
    };
}
