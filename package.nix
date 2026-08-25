{ lib
, stdenv
, buildGoModule
, buildNpmPackage
, jq
, makeWrapper
}:

let
  version = "2026.07.20";

  frontend = buildNpmPackage {
    pname = "jfshare-web";
    inherit version;
    src = ./web;

    npmDepsHash = "sha256-QY5p6XB4XQPbIPwQG+YnCoLf8UyQvxAA89PZsad3yXY=";

    nativeBuildInputs = [ jq ];

    buildPhase = ''
      runHook preBuild
      npm run build
      runHook postBuild
    '';

    installPhase = ''
      runHook preInstall
      mkdir -p $out/share/web
      cp -a dist $out/share/web/dist
      runHook postInstall
    '';
  };

  backend = buildGoModule {
    pname = "jfshare";
    inherit version;
    src = ./.;

    vendorHash = "sha256-9SS7z4Rviv/Ou5CPLx9s71GMvajUaqerYu2d3oxjS2s=";

    subPackages = [ "cmd/server" ];

    nativeBuildInputs = [ jq ];

    doCheck = true;
  };

in
stdenv.mkDerivation {
  pname = "jellyfin-share-backend";
  inherit version;

  nativeBuildInputs = [ makeWrapper ];

  dontUnpack = true;
  dontBuild = true;

  installPhase = ''
    runHook preInstall

    mkdir -p $out/{bin,share/jellyfin-share/web/dist}

    cp ${backend}/bin/server $out/share/jellyfin-share/jfshare
    cp -r ${./migrations} $out/share/jellyfin-share/migrations

    cp -a ${frontend}/share/web/dist/* $out/share/jellyfin-share/web/dist/

    makeWrapper $out/share/jellyfin-share/jfshare $out/bin/jellyfin-share \
      --chdir $out/share/jellyfin-share

    runHook postInstall
  '';

  meta = with lib; {
    description = "Jellyfin share link backend with direct links, audio/subtitle selection";
    license = licenses.mit;
    platforms = platforms.linux;
    mainProgram = "jellyfin-share";
  };
}
