ENTRY  := ./main.go
TARGET := ./api.out
CC     := go
MIGRATION_PATH := ./migrations

# Environment variables for running locally
DB_CONTAINER_NAME="video-meta-db"
DB_USER?="postgres"
DB_PASSWORD?="password"
DB_NAME?="postgres"
DB_HOST?="localhost"
DB_PORT?=5432

all: build

build:
	$(CC) build -o $(TARGET) $(ENTRY)

build-docker:
	docker build \
		--build-arg DB_USER=$(DB_USER) \
		--build-arg DB_NAME=$(DB_NAME) \
		--build-arg DB_PASSWORD=$(DB_PASSWORD) \
		--build-arg DB_HOST=$(DB_HOST) \
		--build-arg DB_PORT=$(DB_PORT) \
		--label "video-meta-service" \
		-f Dockerfile -t video-meta-service .

clean:
	rm -rf $(TARGET)

format:
	$(CC) fmt ./...

local-db-up:
	docker run -d --name $(DB_CONTAINER_NAME) -p 5432:5432 -e POSTGRES_PASSWORD=$(DB_PASSWORD) postgres:15.4
	sleep 3 

local-db-seed:
	PGPASSWORD=$(DB_PASSWORD) psql -d $(DB_NAME) -U $(DB_USER) -h $(DB_HOST) -f ./seed.sql

local-db-connect:
	psql -d $(DB_NAME) -U $(DB_USER) -h $(DB_HOST)

local-db-clean:
	@for i in $$(docker ps -a | awk '/video-meta-db/ {print $$1}'); do docker stop $$i && docker rm $$i; done

# Allows usage of debugger i.e: gdb
debug: clean
	$(CC) build -gcflags=all="-N -l" -ldflags=-compressdwarf=false -o $(TARGET) $(ENTRY)

migrate-up:
	migrate -path $(MIGRATION_PATH) \
		-database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-verbose up

migrate-down:
	migrate -path $(MIGRATION_PATH) \
		-database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-verbose down

migrate-create:
	migrate create -dir $(MIGRATION_PATH) -ext sql $(name)

local-db-init: local-db-clean local-db-up migrate-up local-db-seed 

run: local-db-init
	@MAKE
	@DB_USER=$(DB_USER) \
	PGPASSWORD=$(DB_PASSWORD) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_NAME=$(DB_NAME) \
	DB_HOST=$(DB_HOST) \
	DB_PORT=$(DB_PORT) \
	$(TARGET)
