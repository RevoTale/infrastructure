FROM ghcr.io/biomejs/biome:latest@sha256:bdf4188ce40dc3caa960a893bd64f9d829e59d856241787de6fbb7e2f5ee2740 AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:b8c3669fb4cc652621b3386517c40ab3e151e827de958d6fe75528798ffb5727

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
