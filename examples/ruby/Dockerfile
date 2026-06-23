FROM ruby:3.3-alpine
COPY --from=forte-example-builder /build/bin/forte /usr/local/bin/forte
COPY examples/ruby/app.tar.gz /tmp/ruby-app.tar.gz
RUN adduser -HDSG www-data www-data
RUN forte deploy ruby-app www-data


WORKDIR /srv/ruby-app

RUN gem install rackup webrick

USER www-data
EXPOSE 8000

ENTRYPOINT ["rackup", "--host", "0.0.0.0", "--port", "8000"]
