FROM migrate/migrate

WORKDIR /root/migration

COPY . .

ENTRYPOINT ["/bin/sh", "scripts/migrate.sh"]
