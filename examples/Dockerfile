FROM golang:1.25-alpine AS fortebuilder

ARG FORTEVERSION

WORKDIR /build
COPY cmd/ cmd/
COPY internal/ internal/
COPY go.mod .cz.toml ./
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=${FORTEVERSION}" \
    -o bin/forte \
    ./cmd/forte
