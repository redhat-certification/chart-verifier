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
	go test -v -coverprofile=coverage.out ./...

# Build the container image. Usage: make build-image IMAGE_TAG=my_tag
# If IMAGE_TAG is not provided, use the COMMIT_ID
.PHONY: build-image
build-image:
	# TODO: Adding --no-cache option as a workaround to https://github.com/containers/buildah/issues/4632
	# This can be removed as soon as we can ensure that the ubuntu-latest runner image uses podman=>4.6.0
	$(IMAGE_BUILDER) build \
		--no-cache \
		--label quay.expires-after=$(QUAY_EXPIRE_AFTER) \
		-t $(IMAGE_REPO)/chart-verifier:$(IMAGE_TAG) .

# Push the container image. Usage: make push-image IMAGE_TAG=my_tag
# If IMAGE_TAG is not provided, use the COMMIT_ID
.PHONY: push-image
push-image:
	$(IMAGE_BUILDER) push $(IMAGE_REPO)/chart-verifier:$(IMAGE_TAG)

.PHONY: gosec
gosec: install.gosec
	$(GOSEC) -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir tests ./...

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
	./$(VENV_TOOLS_BIN)/pip install -r requirements.txt
	./$(VENV_TOOLS_BIN)/python setup.py install


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