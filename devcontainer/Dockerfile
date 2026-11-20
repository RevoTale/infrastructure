FROM ghcr.io/biomejs/biome:latest@sha256:0eecaa0c37167e816b75ea0c7b97fef12a9b868ea56ca86e32873ff6eccd87b5 AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:41ff7ce12ec45f2704f27a6e091831186ad922ec8112e6339425bfefdf5a70f9

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
