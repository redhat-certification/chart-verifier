.DEFAULT_GOAL:=help

IMAGE_BUILDER?=podman
IMAGE_REPO?=quay.io/redhat-certification
COMMIT_ID=$(shell git rev-parse --short HEAD)
COMMIT_ID_LONG=$(shell git rev-parse HEAD)
IMAGE_TAG=$(COMMIT_ID)
QUAY_EXPIRE_AFTER="never"

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


# These values capture client-go's supported Kubernetes version and uses that to
# inform some sane defaults for the chart-verifier CLI, particularly when faking server
# interactions for things like template rendering. It's modeled after Helm.
K8S_MODULES_VER=$(subst ., ,$(subst v,,$(shell go list -f '{{.Version}}' -m k8s.io/client-go)))
K8S_MODULES_MAJOR_VER=$(shell echo $$(($(firstword $(K8S_MODULES_VER)) + 1)))
K8S_MODULES_MINOR_VER=$(word 2,$(K8S_MODULES_VER))

LDFLAGS :=
LDFLAGS += -s -w
LDFLAGS += -X github.com/redhat-certification/chart-verifier/cmd.CommitIDLong=$(COMMIT_ID_LONG)
LDFLAGS += -X github.com/redhat-certification/chart-verifier/internal/chartverifier/checks.defaultMockedKubeVersionString=v$(K8S_MODULES_MAJOR_VER).$(K8S_MODULES_MINOR_VER)

.PHONY: bin
bin:
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags '$(LDFLAGS)' \
		-o ./out/chart-verifier main.go

.PHONY: lint
lint: install.golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: bin_win
bin_win:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "-s -w -X 'github.com/redhat-certification/chart-verifier/cmd.CommitIDLong=$(COMMIT_ID_LONG)'" \
		-o .\out\chart-verifier.exe main.go

.PHONY: test
test:
	go test -v -coverprofile=coverage.out ./...

# Build the container image. Usage: make build-image IMAGE_TAG=my_tag
# If IMAGE_TAG is not provided, use the COMMIT_ID
.PHONY: build-image
build-image:
	$(IMAGE_BUILDER) build \
		--label quay.expires-after=$(QUAY_EXPIRE_AFTER) \
		-t $(IMAGE_REPO)/chart-verifier:$(IMAGE_TAG) .

# Push the container image. Usage: make push-image IMAGE_TAG=my_tag
# If IMAGE_TAG is not provided, use the COMMIT_ID
.PHONY: push-image
push-image:
	$(IMAGE_BUILDER) push $(IMAGE_REPO)/chart-verifier:$(IMAGE_TAG)

.PHONY: gosec
gosec: install.gosec
	$(GOSEC) -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir tests --exclude G304 ./...
	# excluding rule G304 because hits currently produce an invalid SARIF.

### Python Specific Targets
PY_BIN ?= python3

# The virtualenv containing code style tools.
VENV_CODESTYLE = venv.codestyle
VENV_CODESTYLE_BIN = $(VENV_CODESTYLE)/bin

# The virtualenv containing our CI scripts
VENV_TOOLS = venv.tools
VENV_TOOLS_BIN = $(VENV_TOOLS)/bin

# This is what we pass to git ls-files.
LS_FILES_INPUT_STR ?= 'scripts/src/*.py' 'tests/*.py'

# The same as format, but will throw a non-zero exit code
# if the formatter had to make changes.
.PHONY: py.ci.format
py.ci.format: py.format
	git diff --exit-code

venv.codestyle:
	$(MAKE) venv.codestyle.always-reinstall

# This target will always install the codestyle venv.
# Useful for development cases.
.PHONY: venv.codestyle.always-reinstall
venv.codestyle.always-reinstall:
	$(PY_BIN) -m venv $(VENV_CODESTYLE)
	./$(VENV_CODESTYLE_BIN)/pip install --upgrade \
		black \
		ruff

.PHONY: py.format
py.format: venv.codestyle
	./$(VENV_CODESTYLE_BIN)/black \
		--verbose \
		$$(git ls-files $(LS_FILES_INPUT_STR))

.PHONY: py.lint
py.lint: venv.codestyle
	./$(VENV_CODESTYLE_BIN)/ruff \
		check \
		$$(git ls-files $(LS_FILES_INPUT_STR))

venv.tools:
	$(MAKE) venv.tools.always-reinstall

# This target will always install the tools at the venv.
# Useful for development cases.
.PHONY: venv.tools.always-reinstall
venv.tools.always-reinstall:
	$(PY_BIN) -m venv $(VENV_TOOLS)
	./$(VENV_TOOLS_BIN)/pip install -r ./scripts/requirements.txt
	./$(VENV_TOOLS_BIN)/pip install ./scripts
	cd ..

gha.lint: actionlint zizmor

actionlint: install.actionlint
	$(ACTIONLINT)

zizmor: install.zizmor
	$(ZIZMOR) .

### Developer Tooling Installation
# gosec
GOSEC = $(shell pwd)/out/gosec
GOSEC_VERSION ?= 6be2b51fd78feca86af91f5186b7964d76cb1256 # v2.22.10
install.gosec: 
	$(call go-install-tool,$(GOSEC),github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION))

# gofumpt
GOFUMPT = $(shell pwd)/out/gofumpt
GOFUMPT_VERSION ?= v0.9.0
install.gofumpt:
	$(call go-install-tool,$(GOFUMPT),mvdan.cc/gofumpt@$(GOFUMPT_VERSION))

# golangci-lint
GOLANGCI_LINT = $(shell pwd)/out/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.2.1
install.golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT):
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))

# actionlint
ACTIONLINT = $(shell pwd)/out/actionlint
ACTIONLINT_VERSION ?= v1.7.0
install.actionlint: $(ACTIONLINT)
$(ACTIONLINT):
	$(call go-install-tool,$(ACTIONLINT),github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION))

# zizmor
ZIZMOR = $(shell pwd)/out/zizmor
ZIZMOR_VERSION ?= v1.17.0
ZIZMOR_ARCH = $(shell uname -m | sed 's/amd64/x86_64/g; s/arm64/aarch64/g')
ZIZMOR_OS = $(shell uname -s | sed 's/Linux/unknown-linux-gnu/g; s/Darwin/apple-darwin/g')
ZIZMOR_PLATFORM ?= $(ZIZMOR_ARCH)-$(ZIZMOR_OS)
install.zizmor: $(ZIZMOR)
$(ZIZMOR):
	$(call github-release-install-tool,$(shell pwd)/out/zizmor.tar.gz,zizmorcore/zizmor,$(ZIZMOR_VERSION),zizmor-$(ZIZMOR_PLATFORM).tar.gz)
	tar -xzf $(shell pwd)/out/zizmor.tar.gz -C $(shell pwd)/out ./zizmor
	rm $(shell pwd)/out/zizmor.tar.gz

# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
GOBIN=$(PROJECT_DIR)/out go install $(2) ;\
}
endef

# github-release-install-tool
# Arguments:
# $1 - destination path
# $2 - GitHub repository
# $3 - release version
# $4 - artifact name
define github-release-install-tool
@[ -f $(1) ] || { \
curl -L https://github.com/$(2)/releases/download/$(3)/$(4) -o $(1) ;\
}
endef