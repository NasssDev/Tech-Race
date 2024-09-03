CURRENT_DIR=$(patsubst %/,%,$(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
ROOT_DIR=$(CURRENT_DIR)
PORT?=8083
RUN_CLOUDINARY=cd pkg/other/cloudinary && go run main/main.go --port=$(PORT)
GO_RUN="go run main/main.go --port=$(PORT)"
AIR=cd pkg/other/cloudinary && PORT=$(PORT) air

run-cloudinarace:
	$(RUN_CLOUDINARY)

run-cloudinary-air:
	$(AIR)

cloudinarace: run-cloudinarace

cloud-air: run-cloudinary-air

