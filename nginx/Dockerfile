FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS routesgen-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/routesgen ./cmd/routesgen

FROM nginx:1.31-alpine

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
