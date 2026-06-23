FROM eclipse-temurin:21-jdk-alpine AS javabuilder
WORKDIR /build
COPY examples/java/src/SimpleServer.java .

RUN javac SimpleServer.java && \
    jar cfe server.jar SimpleServer SimpleServer.class

FROM eclipse-temurin:21-jre-alpine

COPY --from=forte-example-builder /build/bin/forte /usr/local/bin/forte
COPY --from=javabuilder /build/server.jar /tmp/java-app/server.jar
RUN adduser -HDSG www-data www-data
RUN forte deploy java-app www-data

USER www-data
EXPOSE 8000

ENTRYPOINT ["java", "-jar", "/srv/java-app/server.jar"]
