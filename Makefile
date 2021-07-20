Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := -s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)

.PHONY: build
build: esc-build
	go build -ldflags "$(LDFLAGS)" -o build/debug/healthcheck main.go

.PHONY: build-dist
build-dist: esc-build
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o build/release/healthcheck main.go

.PHONY: esc-build
esc-build: build-dashboard
	esc -pkg api -o api/static.go -prefix=dashboard/dist dashboard/dist

.PHONY: run-dashboard
run-dashboard:
	cd dashboard && npm run serve

.PHONY: build-dashboard
build-dashboard:
	cd dashboard && npm run build
