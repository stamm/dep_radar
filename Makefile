# SHELL := /bin/sh
TEST_ARGS?=
MIN_GO_VERSION:=1.9
DEP_VERSION=v0.4.1
DEP_BIN=$(GOPATH)/bin/dep-$(DEP_VERSION)
TMP_DIR=tmp/
COVERAGE_FILE=$(TMP_DIR)coverage.txt
PKGS=$(shell go list -f '{{if len .TestGoFiles}}test-{{.ImportPath}}{{end}}' ./...)

PARALLEL_COUNT?=$(shell getconf _NPROCESSORS_ONLN) #CPU cores

# GOVER:=$(shell go version | cut -f3 -d " " | sed 's/go//')
# IS_DESIRE_VERSION = $(shell expr $(GOVER) \>= $(MIN_GO_VERSION))
# ifeq ($(IS_DESIRE_VERSION),0)
# $(error You have go version $(GOVER), need at least $(MIN_GO_VERSION))
# endif

.PHONY: test
test:
	env GOGC=off go test $(TEST_ARGS) ./...

.PHONY: test
test-race:
	env GOGC=off go test -race $(TEST_ARGS) ./...

.PHONY: generate
generate:
	go generate $(TEST_ARGS) ./...

.PHONY: dep_install
dep_install: $(DEP_BIN)

$(DEP_BIN):
	rm -rf $(GOPATH)/src/github.com/golang/dep/
	mkdir -p $(GOPATH)/src/github.com/golang/dep
	git clone --depth=1 https://github.com/golang/dep.git $(GOPATH)/src/github.com/golang/dep
	cd $(GOPATH)/src/github.com/golang/dep; \
		git pull --tags; \
		git checkout $(DEP_VERSION); \
		go build -o $(DEP_BIN) ./cmd/dep

.PHONY: deps
deps: $(DEP_BIN)
	$(DEP_BIN) ensure -v

.PHONY: bindata
bindata:
	@$(MAKE) -B src/html/templates/bindata.go

src/html/templates/bindata.go: src/html/templates/*.html
	go-bindata -o "./src/html/templates/bindata.go" -ignore "\.go" -pkg "templates" ./src/html/templates/

build: deps src/html/templates/bindata.go
	go build -o $(GOPATH)/bin/dep_radar $(GOPATH)/src/github.com/stamm/dep_radar/cmd/dep_radar/main.go

.PHONY: coverage
coverage: 
	go test -i ./...
	rm -f $(COVERAGE_FILE)
	make -j $(PARALLEL_COUNT) $(PKGS)

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

.PHONY: $(PKGS)
$(PKGS): $(TMP_DIR)
	$(eval $@_package := $(subst test-,,$@))
	$(eval prof_path = $(GOPATH)/src/$($@_package)/profile.out)
	@go test -race -coverprofile=$(prof_path) -covermode=atomic $($@_package); 
	@if [[ -f $(prof_path) ]]; then \
		cat $(prof_path) >> $(COVERAGE_FILE) ;\
		rm $(prof_path) ;\
	fi \

