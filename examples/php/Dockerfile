FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY . .
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2 || echo "0.0.0")" \
    -o bin/forte \
    ./cmd/forte

FROM php:8.5-alpine
COPY --from=builder /build/bin/forte /usr/local/bin/forte
COPY examples/php/app.tar.gz /tmp/php-app.tar.gz
RUN forte deploy php-app www-data

USER www-data
EXPOSE 8000

CMD ["php", "-t", "/srv/php-app", "-S", "0.0.0.0:8000"]
