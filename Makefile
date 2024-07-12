CURRENT_DIR=$(patsubst %/,%,$(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
ROOT_DIR=$(CURRENT_DIR)
DOCKER_COMPOSE?=docker compose
DOCKER_NAME=cloudinarace
DOCKER_EXEC=$(CURRENT_USER) docker exec -it $(DOCKER_NAME) sh
PORT?=8083
RUN_CLOUDINARY=cd pkg/other/cloudinary && go run main/main.go --port=$(PORT)
GO_RUN="go run main/main.go --port=$(PORT)"
AIR=cd pkg/other/cloudinary && air
RUN_DOCKER=cd pkg/other/cloudinary && docker compose up
RUN_TECHRACEAPI=go run cmd/api/main.go

down:
	cd pkg/other/cloudinary && $(DOCKER_COMPOSE) down

build:
	cd pkg/other/cloudinary && $(DOCKER_COMPOSE) up --build --no-recreate -d

up:
	cd pkg/other/cloudinary && $(DOCKER_COMPOSE) up -d

exec:
	$(DOCKER_EXEC)  -c $(GO_RUN)

run-cloudinarace:
	$(RUN_CLOUDINARY)

run-cloudinary-air:
	$(AIR)

run-tech-race:
	$(RUN_TECHRACEAPI)

tidy-cloud:
	cd pkg/other/cloudinary  && go mod tidy

cloud-docker: up exec

cloud-down : down

start-cloudinarace: run-cloudinarace

cloudinarace-tidy: tidy-cloud

cloud-air: run-cloudinary-air

run: run-tech-race