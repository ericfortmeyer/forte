FROM python:3.12-alpine

COPY --from=forte-example-builder /build/bin/forte /usr/local/bin/forte
COPY examples/python/app.tar.gz /tmp/python-app.tar.gz
RUN adduser -HDSG www-data www-data
RUN forte deploy python-app www-data

USER www-data
EXPOSE 8000

ENTRYPOINT ["python3", "/srv/python-app/main.py"]
