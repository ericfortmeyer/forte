FROM golang:1.25-alpine AS fortebuilder

ARG FORTEVERSION

WORKDIR /build
COPY cmd/ cmd/
COPY internal/ internal/
COPY go.mod ./
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=${FORTEVERSION}" \
    -o bin/forte \
    ./cmd/forte

FROM scratch

COPY --from=fortebuilder --chmod=0755 /build/bin/forte /usr/local/bin/forte

ENTRYPOINT ["/usr/local/bin/forte"]
