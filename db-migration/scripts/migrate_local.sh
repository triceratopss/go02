#!/bin/sh
cd `dirname $0`

URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?search_path=public&sslmode=disable"

migrate -source file://../migrations -database ${URL} "$@"
