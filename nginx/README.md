# Nginx Edge Routing

This image now owns TLS termination, shared caching, and host/path routing for the compose stack. Traefik is no longer
part of the runtime path.

## How config is generated

- `nginx/routes.yml` is the canonical route manifest.
- `task gen` runs `go run ./cmd/routesgen` to regenerate `templates/routes/*.conf.template` and the route inventory
  below.
- The image relies on the official `nginx` Docker entrypoint, which renders `/etc/nginx/templates/*.template` into
  `/etc/nginx/conf.d/*.conf` with `envsubst` at container start.
- `docker-compose.base.yml` passes `PROJECT_DOMAIN` into `cache-proxy`, so generated templates remain domain-agnostic. (thats internal info for RevoTale backend)
- Shared cache, TLS, and proxy primitives live in `templates/00-edge.conf.template` and `includes/*.conf`.

## Adding a route

1. Edit `nginx/routes.yml`.
2. Run `task gen`.
3. Review the generated `templates/routes/*.conf.template` diff and the updated route inventory above.
4. Rebuild or recreate `cache-proxy`.
