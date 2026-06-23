FROM ruby:4.0-alpine
COPY --from=forte-example-builder /build/bin/forte /usr/local/bin/forte
COPY examples/ruby/app.tar.gz /tmp/ruby-app.tar.gz
RUN ["gem", "install", "rackup", "webrick"]
RUN ["adduser", "-HDSG", "www-data", "www-data"]
RUN ["forte", "deploy", "ruby-app", "www-data"]

USER www-data
EXPOSE 8000

ENTRYPOINT ["rackup", "--host", "0.0.0.0", "--port", "8000"]
CMD ["--host", "0.0.0.0", "--port", "8000", "/srv/ruby-app/config.ru"]
