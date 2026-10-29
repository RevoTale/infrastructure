FROM ghcr.io/biomejs/biome:latest@sha256:7ca467dc36a0773e8bb92c169a40725c5f5cd0018b034e6b2bb6a991e337cedf AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:71b6fe7c4e07ec747a94a72b2a17f2f8c286eaa7b1c59862888df764efb8b6c6

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
