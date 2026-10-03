FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS routesgen-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/routesgen ./cmd/routesgen

FROM nginx:1.31-alpine@sha256:df221db836e1754089190208cee7eeda94f233197056426eda74a43ab1abeac2

COPY nginx.conf /etc/nginx/nginx.conf
COPY includes/ /etc/nginx/includes/
COPY static/ /usr/share/nginx/html/
COPY docker-entrypoint.d/15-generate-routes.sh /docker-entrypoint.d/15-generate-routes.sh
COPY --from=routesgen-builder /out/routesgen /usr/local/bin/routesgen

RUN chmod +x /docker-entrypoint.d/15-generate-routes.sh /usr/local/bin/routesgen \
    && rm -f /etc/nginx/conf.d/default.conf \
    && mkdir -p /var/cache/nginx /etc/nginx/templates/routes /etc/nginx/conf.d/routes

VOLUME /etc/nginx/certs/
EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
