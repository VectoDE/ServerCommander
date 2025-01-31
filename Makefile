# Makefile für das Build-Management

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
BUILD_DIR := ./bin/$(GOOS)/$(GOARCH)

# Standardziel: Alle Betriebssysteme bauen
build: windows linux macos

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/ServerCommander.exe cmd/server-commander/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/ServerCommander cmd/server-commander/main.go

macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/ServerCommander cmd/server-commander/main.go

clean:
	rm -rf ./bin/*

.PHONY: clean build windows linux macos
