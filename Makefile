.PHONY: build-linux build-mac

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/iago_linux main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o build/iago_mac_amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o build/iago_mac_arm64 main.go
