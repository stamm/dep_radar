# SHELL := /bin/sh
LAST_GOPATH_DIR:=$(lastword $(subst :, ,$(GOPATH)))
GOBIN:=$(LAST_GOPATH_DIR)/bin
TEST_ARGS?=
MIN_GO_VERSION:=1.9
DEP_VERSION=v0.4.1
DEP_BIN=$(GOPATH)/bin/dep-$(DEP_VERSION)
TMP_DIR=tmp/
COVERAGE_FILE=$(TMP_DIR)coverage.txt
PKGS=$(shell go list -f '{{if len .TestGoFiles}}test-{{.ImportPath}}{{end}}' ./...)
RELEASE?=no_version

PARALLEL_COUNT?=$(shell getconf _NPROCESSORS_ONLN) #CPU cores

APP?=dep_radar
CONTAINER_IMAGE?=docker.io/stamm/${APP}
APP_BIN=$(GOPATH)/bin/$(APP)

GOMETALINTER_BIN:=$(GOBIN)/gometalinter.v2
GOLINT_BIN:=$(GOBIN)/golint

GO_BINDATA_BIN=$(GOPATH)/bin/go-bindata


UNAME_S=$(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OS=linux
endif
ifeq ($(UNAME_S),Darwin)
	OS=darwin
endif
UNAME_M=$(shell uname -m)
ARCH=386
ifeq ($(UNAME_M),x86_64)
	ARCH=amd64
endif
# GOVER:=$(shell go version | cut -f3 -d " " | sed 's/go//')
# IS_DESIRE_VERSION = $(shell expr $(GOVER) \>= $(MIN_GO_VERSION))
# ifeq ($(IS_DESIRE_VERSION),0)
# $(error You have go version $(GOVER), need at least $(MIN_GO_VERSION))
# endif

.PHONY: generate
generate:
	go generate $(TEST_ARGS) ./...

.PHONY: dep_install
dep_install: $(DEP_BIN)

$(DEP_BIN):
	curl -L -o $(DEP_BIN) "https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-$(OS)-$(ARCH)"
	chmod +x $(DEP_BIN)

.PHONY: deps
deps: vendor/touch

vendor/touch: $(DEP_BIN) Gopkg.toml Gopkg.lock
	$(DEP_BIN) ensure
	touch $@

.PHONY: bindata
bindata:
	@$(MAKE) -B html/templates/bindata.go

html/templates/bindata.go: html/templates/*.html $(GO_BINDATA_BIN)
	$(GO_BINDATA_BIN) -o "./$@" -ignore "\.go" -pkg "templates" ./html/templates/

$(GO_BINDATA_BIN):
	go get -u github.com/jteeuwen/go-bindata/...

.PHONY: build
build: $(APP_BIN)

$(APP_BIN): vendor/touch html/templates/bindata.go
	go build -ldflags="-s -w" -o $@ $(GOPATH)/src/github.com/stamm/dep_radar/cmd/dep_radar/main.go

run: html/templates/bindata.go
	go run ./cmd/dep_radar/main.go


### RELEASE
release: mkdir_release
	git tag $(RELEASE)
	git push --tags
	@$(MAKE) -B -j3 tmp/release/$(RELEASE)/$(APP)-darwin-amd64.tar.gz tmp/release/$(RELEASE)/$(APP)-linux-amd64.tar.gz docker_latest

build_release: mkdir_release
	@$(MAKE) -B -j3 tmp/release/$(RELEASE)/$(APP)-darwin-amd64.tar.gz tmp/release/$(RELEASE)/$(APP)-linux-amd64.tar.gz

.PHONY: mkdir_release
mkdir_release:
	mkdir -p tmp/release/$(RELEASE)
	rm -rf tmp/release/*

tmp/release/$(RELEASE)/$(APP)-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $@ ./cmd/dep_radar/main.go

tmp/release/$(RELEASE)/$(APP)-darwin-amd64.tar.gz: tmp/release/$(RELEASE)/$(APP)-darwin-amd64
	tar -czf $@ tmp/release/$(RELEASE)/$(APP)-darwin-amd64

tmp/release/$(RELEASE)/$(APP)-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ ./cmd/dep_radar/main.go

tmp/release/$(RELEASE)/$(APP)-linux-amd64.tar.gz: tmp/release/$(RELEASE)/$(APP)-linux-amd64
	tar -czf $@ tmp/release/$(RELEASE)/$(APP)-linux-amd64

### DOCKER IMAGES
.PHONY: docker_build
docker_build:
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

.PHONY: docker_push
docker_push: docker_build
	docker push $(CONTAINER_IMAGE):$(RELEASE)

.PHONY: docker_latest
docker_latest: docker_push
	docker tag $(CONTAINER_IMAGE):$(RELEASE) $(CONTAINER_IMAGE):latest
	docker push $(CONTAINER_IMAGE):latest



### TESTS
.PHONY: test
test: vendor/touch
	env GOGC=off go test $(TEST_ARGS) ./...

.PHONY: test
test-race:
	env GOGC=off CGO_ENABLED=1 go test -race $(TEST_ARGS) ./...

.PHONY: coverage
coverage: $(TMP_DIR) vendor/touch
	env CGO_ENABLED=1 go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...

$(TMP_DIR):
	mkdir -p $(TMP_DIR)



### LINTERS
## install gometalinter
$(GOMETALINTER_BIN):
	go get -u gopkg.in/alecthomas/gometalinter.v2

$(GOLINT_BIN):
	go get -u golang.org/x/lint/golint

lint: $(GOLINT_BIN) $(GOMETALINTER_BIN) | $(TEST_TMP_DIR)
	$(GOMETALINTER_BIN) --no-config --disable-all --enable=golint --deadline=5m -e "html/templates/bindata.go" -e "^vendor/" ./...
