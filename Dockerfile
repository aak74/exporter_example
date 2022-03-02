FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /task-exporter .

##
## Deploy
##
FROM alpine:3.15.0

WORKDIR /

COPY --from=build /task-exporter /task-exporter

EXPOSE 8080

ENTRYPOINT ["/task-exporter"]
