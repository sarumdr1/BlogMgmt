BLOG_SERVICE_BINARY=blogServiceApp
LOGGER_SERVICE_BINARY=loggerServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_blog_service build_logger
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

build_blog_service:
	@echo "Building blog service binary..."
	cd ../blog-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BLOG_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the frond end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## build_front: builds the logger binary
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_SERVICE_BINARY} ./cmd/api
	@echo "Done!"