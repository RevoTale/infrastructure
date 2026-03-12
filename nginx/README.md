# Nginx Edge Routing

This image now owns TLS termination, shared caching, and host/path routing for the compose stack. Traefik is no longer
part of the runtime path.

## Runtime contract

- `PROJECT_DOMAIN` must be present at runtime.
- `/etc/nginx/routes.yml` must be available at runtime, typically via a bind mount.
- `/docker-entrypoint.d/15-generate-routes.sh` deletes stale route templates and stale rendered route config, regenerates
  `/etc/nginx/templates/routes/*.conf.template`, then the official `nginx` entrypoint renders those templates into
  `/etc/nginx/conf.d/routes/*.conf` with `envsubst`.
- `templates/` is runtime-generated route config only.
- Shared cache, HTTP-scope, TLS, and proxy primitives remain static in `includes/http-edge.conf` and `includes/*.conf`.
- Static branded fallback pages for local nginx 4xx/5xx responses live under `static/__error_pages/`.
- Generated route configs keep `proxy_pass http://$upstream;` with `resolver 127.0.0.11` so backends can be redeployed
  without restarting nginx.

## Extending the image

Consumers can either mount a manifest or build a child image:

```dockerfile
FROM ghcr.io/revotale/infrastructure-nginx:latest

COPY routes.yml /etc/nginx/routes.yml
```

## Local verification

- From `infrastructure/nginx`, `task validate` runs the Go tests.
