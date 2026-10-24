FROM ghcr.io/biomejs/biome:latest@sha256:1a4bb81fbb05aaa22fc0a9ba4209db499b942b1d271797d3fe5dcf86f85407a9 AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
