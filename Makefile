.DEFAULT_GOAL:=help

IMAGE_BUILDER?=podman
IMAGE_REPO?=quay.io/redhat-certification
COMMIT_ID=$(shell git rev-parse --short HEAD)
COMMIT_ID_LONG=$(shell git rev-parse HEAD)
IMAGE_TAG=$(COMMIT_ID)

default: bin

.PHONY: all
all:  tidy fmt bin test

# This is a backwards-compatible target to help get
# the PR merged without breaking gha. It can be removed
# once this PR merges.
.PHONY gomod_tidy:
gomod_tidy: tidy

.PHONY gofmt:
gofmt: fmt

.PHONY: tidy
tidy:
	go mod tidy
	git diff --exit-code


.PHONY: fmt
fmt: install.gofumpt
	# -l: list files whose formatting differs from gofumpt's
	# -w: write results to source files instead of stdout
	${GOFUMPT} -l -w . 
	git diff --exit-code

.PHONY: bin
bin:
	CGO_ENABLED=0 go build \
		-ldflags "-X 'github.com/redhat-certification/chart-verifier/cmd.CommitIDLong=$(COMMIT_ID_LONG)'" \
		-o ./out/chart-verifier main.go

.PHONY: lint
lint: install.golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: bin_win
bin_win:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags "-X 'github.com/redhat-certification/chart-verifier/cmd.CommitIDLong=$(COMMIT_ID_LONG)'" \
		-o .\out\chart-verifier.exe main.go

.PHONY: test
test:
	go test -v ./...

# Build the container image. Usage: make build-image IMAGE_TAG=my_tag
# If IMAGE_TAG is not provided, use the COMMIT_ID
.PHONY: build-image
build-image:
	$(IMAGE_BUILDER) build -t $(IMAGE_REPO)/chart-verifier:$(IMAGE_TAG) .

# Push the container image. Usage: make push-image IMAGE_TAG=my_tag
# If IMAGE_TAG is not provided, use the COMMIT_ID
.PHONY: push-image
push-image:
	$(IMAGE_BUILDER) push $(IMAGE_REPO)/chart-verifier:$(IMAGE_TAG) .

.PHONY: gosec
gosec: install.gosec
	$(GOSEC) -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir tests ./...

### Developer Tooling Installation

# gosec
GOSEC = $(shell pwd)/out/gosec
GOSEC_VERSION ?= latest
install.gosec: 
	$(call go-install-tool,$(GOSEC),github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION))

# gofumpt
GOFUMPT = $(shell pwd)/out/gofumpt
install.gofumpt:
	$(call go-install-tool,$(GOFUMPT),mvdan.cc/gofumpt@latest)

# golangci-lint
GOLANGCI_LINT = $(shell pwd)/out/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.52.2
install.golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT):
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))\

# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
GOBIN=$(PROJECT_DIR)/out go install $(2) ;\
}
endef
