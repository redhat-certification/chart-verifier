$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o chart-verifier.exe main.go
