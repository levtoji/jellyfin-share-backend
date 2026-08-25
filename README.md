# Jellyfin Share Backend (Enhanced)

Forked from [monxas/jellyfin-share-backend](https://github.com/monxas/jellyfin-share-backend) with the following additions:

## Features

### Direct Playback Links (`/direct/{token}`)

Generate a single URL that starts playback immediately — no web UI needed. Paste into VRChat, VLC, or any HLS-compatible player.

```
https://watch.example.com/direct/abc123
?AudioStreamIndex=2
&SubtitleStreamIndex=5
&SubtitleMethod=Encode
```

### Audio Track Selection

The share page displays all available audio tracks from the Jellyfin media source. Select a track and the direct link updates with `AudioStreamIndex=<index>`.

### Subtitle Burn-In

Select subtitles from the media's available subtitle tracks. When selected, the direct link includes `SubtitleMethod=Encode` which tells Jellyfin to hardcode the subtitles into the video stream (required for VRChat, which cannot render separate subtitle streams).

### Transcoding Defaults

Pre-configured for quality playback:

| Parameter | Value |
|---|---|
| Video Codec | H.264 |
| Video Bitrate | 4 Mbps |
| Max Resolution | 1080p |
| Allow Stream Copy (video) | No (always transcode) |
| Tone Mapping | Enabled |
| Audio Codec | AAC (for HLS master) |

The HLS stream proxy removes `AudioCodec` from `.ts`/`.m4s` segment URLs to prevent FFmpeg decoding issues.

### Series Support

For Series shares, the backend automatically resolves the first episode's media streams (audio tracks, subtitles) and redirects playback to the correct episode.

## NixOS Module

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    jellyfin-share.url = "github:levtoji/jellyfin-share-backend";
  };

  outputs = { self, nixpkgs, jellyfin-share, ... }: {
    nixosConfigurations.myhost = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      specialArgs = { inherit jellyfin-share; };
      modules = [
        ({ ... }: {
          imports = [ jellyfin-share.nixosModules.default ];

          services.jellyfin-share = {
            enable = true;
            port = 8097;
            databaseUrl = "postgres://watch:password@127.0.0.1:5432/watch?sslmode=disable";
            jellyfinUrl = "http://127.0.0.1:8096";
            jellyfinApiKey = "your-jellyfin-api-key";
            backendApiKey = "your-backend-api-key";
            publicBaseUrl = "https://watch.example.com";
          };
        })
      ];
    };
  };
}
```

### Module Options

| Option | Type | Default | Description |
|---|---|---|---|
| `services.jellyfin-share.enable` | bool | `false` | Enable the service |
| `services.jellyfin-share.user` | string | `jellyfin-share` | Service user |
| `services.jellyfin-share.group` | string | `jellyfin-share` | Service group |
| `services.jellyfin-share.port` | int | `8097` | HTTP listen port |
| `services.jellyfin-share.databaseUrl` | string | — | PostgreSQL DSN |
| `services.jellyfin-share.jellyfinUrl` | string | `http://127.0.0.1:8096` | Jellyfin base URL |
| `services.jellyfin-share.jellyfinApiKey` | string | — | Jellyfin API key |
| `services.jellyfin-share.backendApiKey` | string | — | Backend admin API key |
| `services.jellyfin-share.publicBaseUrl` | string | — | Public URL for share links |
| `services.jellyfin-share.logLevel` | enum | `info` | Log level |
| `services.jellyfin-share.setupPostgresql` | bool | `true` | Auto-create DB + user |
| `services.jellyfin-share.setupJellyfinPlugin` | bool | `true` | Install Jellyfin Share plugin |

## Standalone Package

```bash
nix build github:levtoji/jellyfin-share-backend#default
```

## Upstream Tracking

This is a fork of `monxas/jellyfin-share-backend`. To rebase on new upstream changes:

```bash
git remote add upstream https://github.com/monxas/jellyfin-share-backend.git
git fetch upstream
git rebase upstream/main
```

## License

MIT (same as upstream)
