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
	if [ ! -d $(GOPATH)/src/github.com/golang/dep ] ;\
	then \
		git clone https://github.com/golang/dep.git $(GOPATH)/src/github.com/golang/dep; \
	fi
	cd $(GOPATH)/src/github.com/golang/dep; \
		git pull --tags; \
		git checkout $(DEP_VERSION); \
		go build -o $(DEP_BIN) ./cmd/dep

.PHONY: deps
deps: vendor/touch

vendor/touch: $(DEP_BIN) Gopkg.toml Gopkg.lock
	$(DEP_BIN) ensure
	touch $@

.PHONY: bindata
bindata:
	@$(MAKE) -B src/html/templates/bindata.go

src/html/templates/bindata.go: src/html/templates/*.html
	go-bindata -o "./$@" -ignore "\.go" -pkg "templates" ./src/html/templates/

.PHONY: build
build: $(APP_BIN)

$(APP_BIN): vendor/touch src/html/templates/bindata.go
	go build -ldflags="-s -w" -o $@ $(GOPATH)/src/github.com/stamm/dep_radar/cmd/dep_radar/main.go

### RELEASE
release: mkdir_release
	git tag $(RELEASE)
	git push --tags
	@$(MAKE) -B -j3 tmp/releases/$(RELEASE)/$(APP)-darwin-amd64.tar.gz tmp/releases/$(RELEASE)/$(APP)-linux-amd64.tar.gz docker_latest

build_release: mkdir_release
	@$(MAKE) -B -j3 tmp/release/$(APP)-darwin-amd64.tar.gz tmp/release/$(APP)-linux-amd64.tar.gz

.PHONY: mkdir_release
mkdir_release:
	mkdir -p tmp/release/
	rm -rf tmp/release/*

tmp/release/$(APP)-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $@ ./cmd/dep_radar/main.go

tmp/release/$(APP)-darwin-amd64.tar.gz: tmp/release/$(APP)-darwin-amd64
	tar -czf $@ tmp/release/$(APP)-darwin-amd64

tmp/release/$(APP)-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ ./cmd/dep_radar/main.go

tmp/release/$(APP)-linux-amd64.tar.gz: tmp/release/$(APP)-linux-amd64
	tar -czf $@ tmp/release/$(APP)-linux-amd64

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
	env GOGC=off go test -race $(TEST_ARGS) ./...

.PHONY: coverage
coverage: $(TMP_DIR) vendor/touch
	go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...

$(TMP_DIR):
	mkdir -p $(TMP_DIR)



### LINTERS
## install gometalinter
$(GOMETALINTER_BIN):
	go get -u gopkg.in/alecthomas/gometalinter.v2

$(GOLINT_BIN):
	go get -u golang.org/x/lint/golint

lint: $(GOLINT_BIN) $(GOMETALINTER_BIN) | $(TEST_TMP_DIR)
	$(GOMETALINTER_BIN) --no-config --disable-all --enable=golint --deadline=5m -e "src/html/templates/bindata.go" ./...
