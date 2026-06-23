FROM golang:1.25-alpine AS fortebuilder

WORKDIR /build
COPY cmd/ cmd/
COPY internal/ internal/
COPY go.mod .cz.toml ./
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2 || echo "0.0.0")" \
    -o bin/forte \
    ./cmd/forte
