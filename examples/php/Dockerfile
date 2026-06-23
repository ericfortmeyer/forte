FROM php:8.5-alpine

COPY --from=forte-example-builder /build/bin/forte /usr/local/bin/forte
COPY examples/php/app.tar.gz /tmp/php-app.tar.gz
COPY --chmod=700 --chown=www-data:www-data examples/php/entrypoint.sh /
RUN forte deploy php-app www-data

USER www-data
EXPOSE 8000

ENTRYPOINT ["/entrypoint.sh"]
