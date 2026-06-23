FROM node:24-alpine

COPY --from=forte-example-builder /build/bin/forte /usr/local/bin/forte
COPY examples/nodejs/app.tar.gz /tmp/nodejs-app.tar.gz
RUN forte deploy nodejs-app node

USER node
EXPOSE 8000

ENTRYPOINT ["node", "/srv/nodejs-app/main.js"]
