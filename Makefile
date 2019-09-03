
.PHONY: test docker clean all

COMMIT=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAG_LOCATION=github.com/keikoproj/aws-auth/cmd/cli

LDFLAGS=-ldflags "-X ${LDFLAG_LOCATION}.buildDate=${BUILD} -X ${LDFLAG_LOCATION}.gitCommit=${COMMIT}"

GIT_TAG=$(shell git rev-parse --short HEAD)
IMAGE ?= aws-auth:latest

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/aws-auth github.com/keikoproj/aws-auth
	chmod +x bin/aws-auth

test:
	go test -v ./... -coverprofile coverage.txt
	go tool cover -html=coverage.txt -o coverage.html

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker tag ${IMAGE} eytanavisror/${IMAGE}
	docker push eytanavisror/${IMAGE}