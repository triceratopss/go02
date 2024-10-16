#!/bin/sh
cd `dirname $0`

apk --update --no-cache add curl

curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.1.2/cloud-sql-proxy.linux.amd64
chmod +x cloud-sql-proxy

./cloud-sql-proxy ${INSTANCE_CONNECTION_NAME} & sleep 2;

urlencode() {
    jq -nr --arg v "$1" '$v|@uri'
}

ENCODED_PASSWORD=$(urlencode ${DB_PASSWORD})
URL="postgres://${DB_USER}:${ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?search_path=public&sslmode=disable"

migrate -source file://../migrations -database ${URL} "$@"
