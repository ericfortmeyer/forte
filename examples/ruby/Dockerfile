FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY . .
RUN go build \
    -ldflags "-X github.com/ericfortmeyer/forte/internal/version.version=$(grep -oE 'version = "[^"]+"' .cz.toml | cut -d'"' -f2 || echo "0.0.0")" \
    -o bin/forte \
    ./cmd/forte

FROM ruby:3.3-alpine
COPY --from=builder /build/bin/forte /usr/local/bin/forte
COPY examples/ruby/app.tar.gz /tmp/ruby-app.tar.gz
RUN adduser -HDSG www-data www-data
RUN forte deploy ruby-app www-data


WORKDIR /srv/ruby-app

RUN gem install rackup webrick

USER www-data
EXPOSE 8000

CMD ["rackup", "--host", "0.0.0.0", "--port", "8000"]
