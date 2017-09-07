TEST_ARGS?=
MIN_GO_VERSION:=1.9
DEP_VERSION=v0.3.0
DEP_BIN=$(GOPATH)/bin/dep-$(DEP_VERSION)
TMP_DIR=tmp/
COVERAGE_FILE=$(TMP_DIR)coverage.txt

GOVER:=$(shell go version | cut -f3 -d " " | sed 's/go//')
IS_DESIRE_VERSION = $(shell expr $(GOVER) \>= $(MIN_GO_VERSION))
ifeq ($(IS_DESIRE_VERSION),0)
$(error You have go version $(GOVER), need at least $(MIN_GO_VERSION))
endif

.PHONY: test
test:
	env GOGC=off go test $(TEST_ARGS) ./...

.PHONY: generate
generate:
	go generate $(TEST_ARGS) ./...


$(DEP_BIN):
	rm -rf $(GOPATH)/src/github.com/golang/dep/
	mkdir -p $(GOPATH)/src/github.com/golang/dep
	git clone https://github.com/golang/dep.git $(GOPATH)/src/github.com/golang/dep
	cd $(GOPATH)/src/github.com/golang/dep; \
		git checkout $(DEP_VERSION); \
		go build -o $(DEP_BIN) ./cmd/dep

.PHONY: deps
deps: $(DEP_BIN)
	$(DEP_BIN) ensure -v



.PHONY: coverage
coverage: $(COVERAGE_FILE)

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

# PKGS=$(go list ./... | grep -v /vendor/)
$(COVERAGE_FILE): $(TMP_DIR)
	$(eval PKGS = $(shell go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...))
	for pkg in $(PKGS); do \
		go test -race -coverprofile=$(TMP_DIR)profile.out -covermode=atomic $$pkg; \
		if [[ -f $(TMP_DIR)profile.out ]]; then \
			cat $(TMP_DIR)profile.out >> $(COVERAGE_FILE) ;\
			rm $(TMP_DIR)profile.out ;\
		fi \
	done
