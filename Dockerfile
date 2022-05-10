FROM golang:1.17-alpine as builder
USER root
WORKDIR /app

COPY . /app

RUN mkdir build
RUN apk add git
RUN go mod download; go build -o /main ./

FROM alpine:latest
USER root
LABEL maintainer="spbu-devs"
COPY --from=builder /main /app/build/main

RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 80

WORKDIR /app/build

ENTRYPOINT [ "./main" ]
