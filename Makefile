ENTRY  := ./main.go
TARGET := ./api.out
CC     := go
MIGRATION_PATH := ./database/migrations
TEST_DATABASE := test

# Environment variables for running locally
JWT_SIGNING_KEY?="jwt_signing_key"
DB_USER?="$(shell whoami)"
DB_PASSWORD?=""
DB_NAME?=video_meta
DB_HOST?="localhost"
DB_PORT?=5432

all: build

build:
	$(CC) build -o $(TARGET) $(ENTRY)

clean:
	rm -rf $(TARGET)

local-db-up:
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:15.4

local-db-connect:
	psql -d postgres -U postgres -h localhost

# Allows usage of debugger i.e: gdb
debug: clean
	$(CC) build -gcflags=all="-N -l" -ldflags=-compressdwarf=false -o $(TARGET) $(ENTRY)

format:
	$(CC) fmt ./...

run:
	@DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_NAME=$(DB_NAME) \
	DB_HOST=$(DB_HOST) \
	DB_PORT=$(DB_PORT) \
	JWT_SIGNING_KEY=$(JWT_SIGNING_KEY) \
	$(TARGET)
