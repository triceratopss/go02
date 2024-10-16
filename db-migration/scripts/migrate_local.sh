#!/bin/sh
cd `dirname $0`

urlencode() {
    jq -nr --arg v "$1" '$v|@uri'
}

ENCODED_PASSWORD=$(urlencode ${DB_PASSWORD})
URL="postgres://${DB_USER}:${ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?search_path=public&sslmode=disable"

migrate -source file://../migrations -database ${URL} "$@"
