CURRENT_DIR=$(patsubst %/,%,$(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
ROOT_DIR=$(CURRENT_DIR)
DOCKER_COMPOSE?=docker-compose
PORT?=8083
RUN_CLOUDINARY=cd pkg/other/cloudinary && go run main/main.go --port=$(PORT)
GO_RUN=go run .
AIR=cd pkg/other/cloudinary && air

down:
	cd pkg/other/cloudinary && $(DOCKER_COMPOSE) down

build:
	cd pkg/other/cloudinary && $(DOCKER_COMPOSE) up --build --no-recreate -d

up:
	cd pkg/other/cloudinary && $(DOCKER_COMPOSE) up -d

run-cloudinarace:
	$(RUN_CLOUDINARY)

run-cloudinary-air:
	$(AIR)



cloud-docker: up

shut-docker : down

start-cloudinarace: run-cloudinarace

cloud-air: run-cloudinary-air