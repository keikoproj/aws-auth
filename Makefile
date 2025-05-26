
.PHONY: test docker clean all

COMMIT=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAG_LOCATION=github.com/keikoproj/aws-auth/cmd/cli

LDFLAGS=-ldflags "-X ${LDFLAG_LOCATION}.buildDate=${BUILD} -X ${LDFLAG_LOCATION}.gitCommit=${COMMIT}"

GIT_TAG=$(shell git rev-parse --short HEAD)
IMAGE ?= aws-auth:latest

all: lint test build

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/aws-auth github.com/keikoproj/aws-auth
	chmod +x bin/aws-auth

test: fmt vet
	go test -v ./... -coverprofile coverage.txt
	go tool cover -html=coverage.txt -o coverage.html

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker push ${IMAGE}

LOCALBIN = $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

GOLANGCI_LINT_VERSION := v2.1.1
GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
.PHONY: golangci-lint
$(GOLANGCI_LINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo "Running golangci-lint"
	$(GOLANGCI_LINT) run ./...

.PHONY: clean
clean:
	@rm -rf ./bin
