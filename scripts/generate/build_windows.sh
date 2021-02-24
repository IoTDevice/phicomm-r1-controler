export CGO_ENABLED=0
export GO111MODULE=on
export GOARCH=386
export GOOS=windows
go build -ldflags -w ./

export GOARCH=amd64
export GOOS=darwin