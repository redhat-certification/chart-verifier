$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o .\out\chart-verifier.exe main.go
