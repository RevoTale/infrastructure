FROM ghcr.io/biomejs/biome:latest@sha256:7cab88dee7200bf88b087dd739b29bbb56dc92d33f7a03978436294b6aca897a AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:71b6fe7c4e07ec747a94a72b2a17f2f8c286eaa7b1c59862888df764efb8b6c6

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
