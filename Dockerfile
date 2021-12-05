FROM node:14-alpine as angular-build-stage

COPY ./Frontend /app/frontend/src

WORKDIR /app/frontend/src

RUN npm install && npx ng build --prod

# Switch to maven build image
FROM maven:3-openjdk-8 as maven-build-stage

# Copy the compiled WebApp from stage 1 to the new container for use as templates
COPY --from=angular-build-stage /app/frontend/src/dist/UrlShortener /app/frontend/dist/

COPY ./Backend /app/backend

WORKDIR /app/backend

RUN mvn install && \
    mvn package

# This is were it gets run
FROM openjdk:8-jre

COPY --from=maven-build-stage /app/backend/target/Url-Shortener-1.0.jar /app/backend/
COPY --from=maven-build-stage /app/frontend /app/frontend

ENV ANGULAR_FRONTEND_DIR=/app/frontend/dist/

# The server port is hardcoded because it doesn't need to be changed here
# If you want it to be different change the port mapping in your docker run command or docker-compose file.
ENV SERVER_PORT=8999

CMD [ "java", "-jar", "/app/backend/Url-Shortener-1.0.jar" ]
