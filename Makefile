
default: bin

.PHONY: all
all:  gomod_tidy gofmt bin test

.PHONY: gomod_tidy
gomod_tidy:
	go mod tidy

.PHONY: gofmt
gofmt:
	go fmt -x ./...

.PHONY: bin
bin:
	 go build -o ./out/chart-verifier main.go

.PHONY: bin_win
bin_win:
	env GOOS=windows GOARCH=amd64 go build -o .\out\chart-verifier.exe main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: build-image
build-image:
	hack/build-image.sh
