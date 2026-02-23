.PHONY: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hdinfo-linux-amd64 ./cmd/hdinfo
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o hdinfo-linux-arm64 ./cmd/hdinfo
