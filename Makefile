BIN_DIR               = bin
BINS                  = bin/notify_slack
GOLANGCI_LINT         = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION = v1.46.0
V                     = 0

RM                    = rm
SHELL                 = /usr/bin/env bash -o pipefail

GO                    = go
GOBUILD               = $(GO) build
GOGET                 = $(GO) get
GOPATH                = $(shell $(GO) env GOPATH)
CGO                   = CGO_ENABLED=0 GOOS=linux
LDFLAGS               = -ldflags="-s -w"
GO_SRC                = $(shell find ./ -name '*.go')
UPPER                 = $(shell echo '$1' | tr '[:lower:]' '[:upper:]')
FORMAT                = sed "s/^/    /g"

ECHO              = @echo "  "
UNAME_S=$(shell uname -s)
ifeq ($(UNAME_S),Linux)
	ECHO=@echo -e "  "
endif

ifneq ($(V),1)
	Q = @
	DEBUG = 2>/dev/null
endif

.PHONY: all
all: help

##@ General
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development
.PHONY: build
build: $(BINS) ## Build executables
$(BIN_DIR)/%: $(GO_SRC)
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(CGO) $(GOBUILD) $(LDFLAGS) -o $@ main.go 2>&1 | $(FORMAT)

.PHONY: fmt
fmt: ## Run go fmt against code
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(GO) fmt ./... 2>&1 | $(FORMAT)

$(GOLANGCI_LINT):
	$(Q)$(ECHO) "GOLANGCI_LINT"
	$(Q) wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh $(DEBUG) | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) $(DEBUG) | $(FORMAT)

.PHONY: lint
lint: ${GOLANGCI_LINT} ## Run golangci-lint linter
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(GOLANGCI_LINT) run --color always 2>&1 | $(FORMAT)

.PHONY: clean
clean: ## Remove bin dir
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(RM) -rf bin/*
