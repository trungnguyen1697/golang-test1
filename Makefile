.PHONY: clean test security build run docker-up docker-down migrate swagger

APP_NAME = golang-test1
BUILD_DIR = ./build
MIGRATIONS_FOLDER = ./platform/migrations
DB_NAME = ecommerce
DB_USER = postgres
DB_PASS = postgres
DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@localhost/$(DB_NAME)?sslmode=disable

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf *.out

swag:
	swag init

build: swag clean
	CGO_ENABLED=0  go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

docker.run: docker.setup docker.postgres docker.fiber migrate.down migrate.up
	@echo "\n===========FGB==========="
	@echo "App is running...\nVisit: http://localhost:5000 OR http://localhost:5000/swagger/"

docker.setup:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network
	docker volume create fibergb-pgdata

docker.fiber.build: swag
	docker build -t golang-test1:latest .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name trungnt-api \
		--network dev-network \
		-p 5000:5000 \
		golang-test1

docker.fiber.run.dev:
	docker run --rm -d \
		--name trungnt-api \
		--network dev-network \
		-p 5000:5000 \
		golang-test1

docker.postgres:
	docker run --rm -d \
		--name trungnt-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=ecommerce \
		-v fibergb-pgdata:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.stop: docker.stop.fiber docker.stop.postgres

docker.stop.fiber:
	docker stop trungnt-api || true

docker.stop.postgres:
	docker stop trungnt-postgres || true

docker.dev:
	docker-compose up