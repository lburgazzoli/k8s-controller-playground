
PROJECT_NAME ?= playgound-operator
PROJECT_VERSION ?= 0.0.1

CONTAINER_REGISTRY ?= quay.io
CONTAINER_REGISTRY_ORG ?= lburgazzoli
CONTAINER_IMAGE_VERSION ?= $(PROJECT_VERSION)
CONTAINER_IMAGE ?= $(CONTAINER_REGISTRY)/$(CONTAINER_REGISTRY_ORG)/$(PROJECT_NAME):$(CONTAINER_IMAGE_VERSION)

LINT_GOGC ?= 10
LINT_TIMEOUT ?= 10m

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCALBIN := $(PROJECT_PATH)/bin

## Tool Versions
CODEGEN_VERSION ?= v0.30.3
KUSTOMIZE_VERSION ?= v5.4.2
CONTROLLER_TOOLS_VERSION ?= v0.16.0
KIND_VERSION ?= v0.23.0
KIND_IMAGE_VERSION ?= v1.30.4
LINTER_VERSION ?= v1.61.0
GOVULNCHECK_VERSION ?= latest
KO_VERSION ?= latest

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
LINTER ?= $(LOCALBIN)/golangci-lint
GOIMPORT ?= $(LOCALBIN)/goimports
YQ ?= $(LOCALBIN)/yq
KIND ?= $(LOCALBIN)/kind
GOVULNCHECK ?= $(LOCALBIN)/govulncheck
KO ?= $(LOCALBIN)/ko
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build


ifndef ignore-not-found
  ignore-not-found = false
endif

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-\/]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: codegen-tools-install ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(PROJECT_PATH)/hack/scripts/gen_crd.sh $(PROJECT_PATH)

.PHONY: generate
generate: codegen-tools-install ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(PROJECT_PATH)/hack/scripts/gen_res.sh $(PROJECT_PATH)
	$(PROJECT_PATH)/hack/scripts/gen_client.sh $(PROJECT_PATH)

.PHONY: fmt
fmt: goimport ## Run go fmt, gomiport against code.
	$(GOIMPORT) -l -w .
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -ldflags="$(GOLDFLAGS)" -o bin/$(PROJECT_NAME) cmd/main.go

.PHONY: run
run: ## Run a controller from your host.
	go run -ldflags="$(GOLDFLAGS)" cmd/main.go run \
		--leader-election=false \
		--zap-devel \
		--health-probe-bind-address ":0" \
		--metrics-bind-address ":0"

.PHONY: run/local
run/local: install ## Install and Run a controller from your host.
	go run -ldflags="$(GOLDFLAGS)" cmd/main.go run \
		--leader-election=false \
		--zap-devel \
		--health-probe-bind-address ":0" \
		--metrics-bind-address ":0"


.PHONY: deps
deps:  ## Tidy up deps.
	go mod tidy


.PHONY: check
check: check/lint  check/vuln

.PHONY: check/lint
check/lint: golangci-lint
	@echo "run golangci-lint"
	@$(LINTER) run \
		--config .golangci.yml \
		--out-format tab \
		--exclude-dirs etc \
		--timeout $(LINT_TIMEOUT) \
		--verbose

.PHONY: check/lint/fix
check/lint/fix: golangci-lint
	@$(LINTER) run \
		--config .golangci.yml \
		--out-format tab \
		--exclude-dirs etc \
		--timeout $(LINT_TIMEOUT) \
		--fix

.PHONY: check/vuln
check/vuln: govulncheck
	@echo "run govulncheck"
	@$(GOVULNCHECK) ./...

##@ Deployment

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

## Location to install dependencies to
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || \
	GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)


.PHONY: golangci-lint
golangci-lint: $(LINTER)
$(LINTER): $(LOCALBIN)
	@test -s $(LOCALBIN)/golangci-lint || \
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINTER_VERSION)

.PHONY: goimport
goimport: $(GOIMPORT)
$(GOIMPORT): $(LOCALBIN)
	@test -s $(LOCALBIN)/goimport || \
	GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@latest

.PHONY: yq
yq: $(YQ)
$(YQ): $(LOCALBIN)
	@test -s $(LOCALBIN)/yq || \
	GOBIN=$(LOCALBIN) go install github.com/mikefarah/yq/v4@latest

.PHONY: kind
kind: $(KIND)
$(KIND): $(LOCALBIN)
	@test -s $(LOCALBIN)/kind || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/kind@$(KIND_VERSION)

.PHONY: ko
ko: $(KO)
$(KO): $(LOCALBIN)
	@test -s $(LOCALBIN)/ko || \
	GOBIN=$(LOCALBIN) go install github.com/google/ko@$(KO_VERSION)

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK)
$(GOVULNCHECK): $(LOCALBIN)
	@test -s $(GOVULNCHECK) || \
	GOBIN=$(LOCALBIN) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

.PHONY: codegen-tools-install
codegen-tools-install: $(LOCALBIN)
	@echo "Installing code gen tools"
	$(PROJECT_PATH)/hack/scripts/install_gen_tools.sh $(PROJECT_PATH) $(CODEGEN_VERSION) $(CONTROLLER_TOOLS_VERSION)
