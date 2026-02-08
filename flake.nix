{
  description = "URL Shortener Nix flake with systemd template module";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  };

  outputs = { self, nixpkgs }:
    let
      lib = nixpkgs.lib;
    in {
      nixosModules.url-shortener = { config, lib, pkgs, ... }:
        let
          cfg = config.services.url-shortener;
        in {
          options.services.url-shortener = {
            enable = lib.mkEnableOption "URL Shortener service";

            jarPath = lib.mkOption {
              type = lib.types.path;
              description = "Path to the Url-Shortener-1.0.jar file.";
            };

            frontendDir = lib.mkOption {
              type = lib.types.path;
              default = "/var/lib/url-shortener/frontend";
              description = "Path to the static frontend assets directory.";
            };

            javaPackage = lib.mkOption {
              type = lib.types.package;
              default = pkgs.temurin-bin-21;
              description = "Java runtime package used to run the service.";
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
                ExecStart = "${cfg.javaPackage}/bin/java -jar ${cfg.jarPath}";
                Restart = "on-failure";
              };

              environment = {
                FRONTEND_DIR = cfg.frontendDir;
                SERVER_PORT = "%i";
              };
            };
          };
        };
    };
}
