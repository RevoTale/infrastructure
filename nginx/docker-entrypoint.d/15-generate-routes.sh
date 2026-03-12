#!/bin/sh
set -eu

config_path="${NGINX_ROUTES_CONFIG_PATH:-/etc/nginx/routes.yml}"
templates_dir="${NGINX_ROUTES_TEMPLATES_DIR:-/etc/nginx/templates/routes}"
rendered_dir="${NGINX_ROUTES_RENDERED_DIR:-/etc/nginx/conf.d/routes}"
routesgen_bin="${NGINX_ROUTESGEN_BIN:-/usr/local/bin/routesgen}"

if [ ! -f "$config_path" ]; then
    echo >&2 "nginx-routes: missing routes manifest at $config_path"
    exit 1
fi

rm -rf "$templates_dir" "$rendered_dir"
mkdir -p "$templates_dir" "$rendered_dir"

"$routesgen_bin" --config "$config_path" --output-dir "$templates_dir" --clean
