FROM golang:1.13-alpine

RUN apk update && apk add --no-cache git && apk add wget && apk add --update -netcat-openbsd

ENV REPO_URL=github.com/sauravgsh16/api-doorway

ENV GOPATH=/app

ENV APP_PATH=${GOPATH}/&{REPO_URL}

ENV WORK_PATH=&{APP_PATH}

LABEL maintainer="Saurav Ghosh <sauravgsh16@gmail.com>"

WORKDIR ${WORK_PATH}

COPY . .

RUN go mod download

RUN go build -o doorway ./cmd

EXPOSE 5000

RUN wget https://raw.githubusercontent.com/eficode/wait-for/master/wait-for

RUN chmod +x wait-for

CMD ["sh", "wait-for", "db:5432", "--", "./doorway"]