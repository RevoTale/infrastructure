FROM ghcr.io/biomejs/biome:latest@sha256:bdf4188ce40dc3caa960a893bd64f9d829e59d856241787de6fbb7e2f5ee2740 AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:1f851004adcd3dff3776b4a1da86727cf52280b0fdc9d5b0568c9c7d74274286

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
