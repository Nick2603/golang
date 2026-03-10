# Build stage
FROM golang:1.25.4-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal/ internal/
COPY cmd/ cmd/

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /tmp/mongodb-http-server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /tmp/mongodb-http-server /app/server

EXPOSE 8080

CMD ["/app/server"]
