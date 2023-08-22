# Multi stage build, do the building in the builder container
FROM golang:1.21-0-alpine3.17 AS builder

ARG DB_USER
ARG DB_PASSWORD
ARG DB_PORT
ARG DB_NAME
ARG DB_HOST

ENV APP api.out
ENV APP_DIR $GOPATH/src/github.com/crypto-server

WORKDIR $APP_DIR

RUN apk add make

COPY ./src ./src
COPY main.go .
COPY go.mod .
COPY go.sum .
COPY Makefile .

RUN go mod download
RUN make

# This is what would get sent to ECR
FROM golang:1.21-0-alpine3.17

ARG DB_USER
ARG DB_PASSWORD
ARG DB_PORT
ARG DB_NAME
ARG DB_HOST

ENV DB_USER $DB_USER
ENV DB_PASSWORD $DB_PASSWORD
ENV DB_PORT $DB_PORT
ENV DB_NAME $DB_NAME
ENV DB_HOST $DB_HOST

WORKDIR $APP_DIR
COPY --chown=0:0 --from=builder $GOPATH/src/github.com/crypto-server/api.out ./api.out
EXPOSE 8080

CMD ["./api.out"]
