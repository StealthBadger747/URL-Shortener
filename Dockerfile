FROM golang:1.22-alpine AS build

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY static ./static

RUN go build -o /app/url-shortener ./cmd/url-shortener

FROM alpine:3.19
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app
COPY --from=build /app/url-shortener /app/url-shortener
COPY --from=build /app/static /app/static

ENV FRONTEND_DIR=/app/static
EXPOSE 8080

CMD ["/app/url-shortener"]
