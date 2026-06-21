FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY . .
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2 || echo "0.0.0")" \
    -o bin/forte \
    ./cmd/forte

FROM python:3.12-slim

COPY --from=builder /build/bin/forte /usr/local/bin/forte
COPY examples/python/app.tar.gz /tmp/python-app.tar.gz
RUN forte deploy python-app www-data

USER www-data
EXPOSE 8000

CMD ["python3", "/srv/python-app/main.py"]
