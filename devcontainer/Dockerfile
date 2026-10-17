FROM ghcr.io/biomejs/biome:latest@sha256:e73b0563691da8c60ebfdb1553adea0c5d531752460e1672c1020d2d83a7e5d7 AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
