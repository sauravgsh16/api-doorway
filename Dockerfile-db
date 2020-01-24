FROM postgres:alpine

RUN apk update && apk upgrade && apk add bash

ADD ./dbinit /docker-entrypoint-initdb.d/