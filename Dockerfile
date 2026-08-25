FROM scratch

WORKDIR /app

COPY dist/deutsch-linux-amd64 /app/deutsch
COPY dist/deutsch-dbinit-linux-amd64 /app/deutsch-dbinit
COPY assets /app/assets
COPY scripts/schema/init_data.sql /app/scripts/schema/init_data.sql

EXPOSE 8888

ENTRYPOINT ["/app/deutsch"]
CMD ["-f", "/app/etc/deutsch.yaml"]
