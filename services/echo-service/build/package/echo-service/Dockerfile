FROM golang:1-alpine AS build-env

COPY . /project

RUN set -xe \
    && cd /project/cmd/echo-service \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o echo-service

FROM scratch

COPY --from=build-env /project/cmd/echo-service/echo-service /app/

EXPOSE 8000

ENTRYPOINT ["/app/echo-service"]
