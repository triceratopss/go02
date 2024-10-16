FROM migrate/migrate

RUN apk add --no-cache jq

WORKDIR /root/migration

COPY . .

ENTRYPOINT ["/bin/sh", "scripts/migrate.sh"]
