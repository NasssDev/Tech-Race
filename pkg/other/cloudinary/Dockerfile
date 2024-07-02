# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22.3 AS Builder
WORKDIR /app

COPY . .
RUN go mod download \
    && go mod verify

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o /build/cloudinarace ./main/main.go

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /appgo

COPY --from=Builder /build/cloudinarace ./cloudinarace

ENTRYPOINT ["./cloudinarace"]