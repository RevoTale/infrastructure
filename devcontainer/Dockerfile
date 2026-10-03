FROM ghcr.io/biomejs/biome:latest@sha256:e73b0563691da8c60ebfdb1553adea0c5d531752460e1672c1020d2d83a7e5d7 AS biome

FROM mcr.microsoft.com/vscode/devcontainers/base@sha256:1f851004adcd3dff3776b4a1da86727cf52280b0fdc9d5b0568c9c7d74274286

COPY --from=biome /usr/local/bin/biome /usr/local/bin/biome
