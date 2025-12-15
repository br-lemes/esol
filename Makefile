.PHONY: build linux release test version windows

build: test
	@go build -ldflags "-s -w"

linux: test
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"

release: version linux windows
	@go run ./tools/release/main.go

test:
	@go test ./...

version: test
	@go run ./tools/version/main.go

windows: test
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"
