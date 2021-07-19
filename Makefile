Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := -s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o build/debug/healthcheck main.go

.PHONY: build-dist
build-dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o build/release/healthcheck main.go
