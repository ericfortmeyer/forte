FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY . .
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2 || echo "0.0.0")" \
    -o bin/forte \
    ./cmd/forte

FROM node:24-alpine

COPY --from=builder /build/bin/forte /usr/local/bin/forte
COPY examples/nodejs/app.tar.gz /tmp/nodejs-app.tar.gz
RUN forte deploy nodejs-app node

USER node
EXPOSE 8000

CMD ["node", "/srv/nodejs-app/main.js"]
