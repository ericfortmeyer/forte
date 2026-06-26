FROM debian:bookworm-slim

RUN ["apt-get", "update"]
RUN ["apt-get", "install", "-y", "bats"]

WORKDIR /usr/local/src

ENTRYPOINT ["bats"]
CMD ["-rx", "./integration_tests/"]
