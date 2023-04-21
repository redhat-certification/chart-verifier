
default: bin

.PHONY: all
all:  gomod_tidy gofmt bin test

.PHONY: gomod_tidy
gomod_tidy:
	go mod tidy

.PHONY: gofmt
gofmt:
	go fmt -x ./...

.PHONY: fmt
fmt: install.gofumpt
	# -l: list files whose formatting differs from gofumpt's
	# -w: write results to source files instead of stdout
	${GOFUMPT} -l -w . 
	git diff --exit-code

.PHONY: bin
bin:
	 go build -o ./out/chart-verifier main.go

.PHONY: lint
lint: install.golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: bin_win
bin_win:
	env GOOS=windows GOARCH=amd64 go build -o .\out\chart-verifier.exe main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: build-image
build-image:
	hack/build-image.sh

.PHONY: gosec
gosec: install.gosec
	$(GOSEC) -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir tests ./...

# Developer Tooling Installation
GOSEC = $(shell pwd)/out/gosec
GOSEC_VERSION ?= latest
install.gosec: 
	$(call go-install-tool,$(GOSEC),github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION))

GOFUMPT = $(shell pwd)/out/gofumpt
install.gofumpt:
	$(call go-install-tool,$(GOFUMPT),mvdan.cc/gofumpt@latest)

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