FROM maven:3-eclipse-temurin-21-alpine as maven-build-stage

COPY ./FrontendHtmx /app/frontend

COPY ./Backend /app/backend

WORKDIR /app/backend

RUN mvn install && \
    mvn package

# This is were it gets run
FROM eclipse-temurin:21-alpine

COPY --from=maven-build-stage /app/backend/target/Url-Shortener-1.0.jar /app/backend/
COPY --from=maven-build-stage /app/frontend /app/frontend

ENV FRONTEND_DIR=/app/frontend

# The server port is hardcoded because it doesn't need to be changed here
# If you want it to be different change the port mapping in your docker run command or docker-compose file.
ENV SERVER_PORT=8999

CMD [ "java", "-jar", "/app/backend/Url-Shortener-1.0.jar" ]
