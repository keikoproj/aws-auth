
.PHONY: test docker clean all

COMMIT=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAG_LOCATION=github.com/eytan-avisror/aws-auth/cmd/cli

LDFLAGS=-ldflags "-X ${LDFLAG_LOCATION}.buildDate=${BUILD} -X ${LDFLAG_LOCATION}.gitCommit=${COMMIT}"

GIT_TAG=$(shell git rev-parse --short HEAD)

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/aws-auth github.com/eytan-avisror/aws-auth/cmd
	chmod +x bin/eks-bootstrapper
