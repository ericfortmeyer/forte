FROM ubuntu:24.04

ENV TERM=xterm

RUN apt-get update && apt-get install -y bats

WORKDIR /usr/local/src
