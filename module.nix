{ config, lib, pkgs, ... }:

let
  cfg = config.services.jellyfin-share;

  jellyfin-share-backend = pkgs.callPackage ./package.nix {};

  jellyfinSharePlugin = pkgs.fetchurl {
    url = "https://github.com/monxas/jellyfin-share-plugin/releases/download/v1.2.1/Jellyfin.Plugin.Share.dll";
    hash = "sha256-LyPm+rYLfNl0RAXayXQF6UKdjhWXsmDjm6oviD0Galw=";
  };
in {
  options.services.jellyfin-share = {
    enable = lib.mkEnableOption "Jellyfin Share backend";

    user = lib.mkOption {
      type = lib.types.str;
      default = "jellyfin-share";
      description = "User to run the service as";
    };

    group = lib.mkOption {
      type = lib.types.str;
      default = "jellyfin-share";
      description = "Group for the service user";
    };

    port = lib.mkOption {
      type = lib.types.port;
      default = 8097;
      description = "Port for the share backend HTTP server";
    };

    databaseUrl = lib.mkOption {
      type = lib.types.str;
      example = "postgres://user:pass@127.0.0.1:5432/watch?sslmode=disable";
      description = "PostgreSQL connection string";
    };

    jellyfinUrl = lib.mkOption {
      type = lib.types.str;
      default = "http://127.0.0.1:8096";
      description = "Jellyfin server base URL";
    };

    jellyfinApiKey = lib.mkOption {
      type = lib.types.str;
      description = "Jellyfin API key";
    };

    backendApiKey = lib.mkOption {
      type = lib.types.str;
      description = "Backend API key for admin access";
    };

    publicBaseUrl = lib.mkOption {
      type = lib.types.str;
      example = "https://watch.example.com";
      description = "Public-facing base URL for share links";
    };

    logLevel = lib.mkOption {
      type = lib.types.enum [ "debug" "info" "warn" "error" ];
      default = "info";
      description = "Log level";
    };

    setupPostgresql = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = "Automatically create PostgreSQL database and user";
    };

    setupJellyfinPlugin = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = "Install the Jellyfin Share plugin and inject loader script into jellyfin-web";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.jellyfin-share = {
      description = "Jellyfin Share Backend";
      after = [ "network.target" ] ++ lib.optionals cfg.setupPostgresql [ "postgresql.service" ];
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        ExecStart = lib.getExe jellyfin-share-backend;
        WorkingDirectory = "${jellyfin-share-backend}/share/jellyfin-share";

        User = cfg.user;
        Group = cfg.group;

        Restart = "on-failure";
        RestartSec = "5s";

        NoNewPrivileges = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        ReadWritePaths = [ "/var/lib/jellyfin-share" ];
        PrivateTmp = true;
      };

      environment = {
        JFSHARE_PORT = toString cfg.port;
        JFSHARE_DB_DSN = cfg.databaseUrl;
        JFSHARE_JELLYFIN_BASE_URL = cfg.jellyfinUrl;
        JFSHARE_JELLYFIN_API_KEY = cfg.jellyfinApiKey;
        JFSHARE_BACKEND_API_KEY = cfg.backendApiKey;
        JFSHARE_PUBLIC_BASE_URL = cfg.publicBaseUrl;
        JFSHARE_LOG_LEVEL = cfg.logLevel;
      };
    };

    systemd.tmpfiles.rules = [
      "d /var/lib/jellyfin-share 0750 ${cfg.user} ${cfg.group} -"
    ] ++ lib.optionals cfg.setupJellyfinPlugin [
      "d /var/lib/jellyfin/plugins/JellyfinShare 0750 jellyfin jellyfin -"
      "C+ /var/lib/jellyfin/plugins/JellyfinShare/Jellyfin.Plugin.Share.dll 0644 jellyfin jellyfin - ${jellyfinSharePlugin}"
    ];

    users.users.${cfg.user} = {
      isSystemUser = true;
      group = cfg.group;
    };
    users.groups.${cfg.group} = {};

    services.postgresql = lib.mkIf cfg.setupPostgresql {
      ensureDatabases = [ "watch" ];
      ensureUsers = [{
        name = "watch";
        ensureDBOwnership = true;
      }];
    };

    nixpkgs.overlays = lib.mkIf cfg.setupJellyfinPlugin [
      (final: prev: {
        jellyfin-web = prev.jellyfin-web.overrideAttrs (oldAttrs: {
          postInstall = (oldAttrs.postInstall or "") + ''
            substituteInPlace $out/share/jellyfin-web/index.html \
              --replace-fail '</body>' '<script src="/plugins/share/loader.js"></script></body>'
          '';
        });
      })
    ];
  };
}
