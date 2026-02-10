FROM ghcr.io/biomejs/biome:latest AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
