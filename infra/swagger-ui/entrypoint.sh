#!/bin/sh

set -e

: "${BASE_URL:=/}"
BASE_URL=$(echo "$BASE_URL" | sed 's:/*$::')

TARGET=/usr/share/nginx/html/swagger
mkdir -p "$TARGET"

cp /service-auth-docs/swagger.json /usr/share/nginx/html/swagger/service-auth.json
cp /service-products-docs/swagger.json /usr/share/nginx/html/swagger/service-products.json

cat <<EOF > /usr/share/nginx/html/swagger/swagger-config.yaml
urls:
  - url: "/api/swagger/swagger/service-auth.json"
    name: "Service AUTH"
  - url: "/api/swagger/swagger/service-products.json"
    name: "Service PRODUCTS"
EOF

export CONFIG_URL=/api/swagger/swagger/swagger-config.yaml
export BASE_URL=/api/swagger

exec /docker-entrypoint.sh nginx -g 'daemon off;' > /dev/null 2>&1