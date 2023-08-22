# Multi stage build, do the building in the builder container
FROM golang:alpine3.18 AS builder

ARG DB_USER
ARG DB_PASSWORD
ARG DB_PORT
ARG DB_NAME
ARG DB_HOST

ENV APP api.out
ENV APP_DIR /app/

WORKDIR $APP_DIR

RUN apk add make

COPY ./server ./server/
COPY main.go .
COPY go.mod .
COPY go.sum .
COPY Makefile .

RUN go mod download
RUN make

# This is what would get sent to ECR
FROM golang:alpine3.18

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

# Copy binary from builder stage
COPY --from=builder /app/api.out ./api.out

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api.out"]
