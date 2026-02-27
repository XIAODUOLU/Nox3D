.PHONY: build clean test run install

# Binary name
BINARY=nox

# Build the project
build:
	go build -o ./bin/$(BINARY) ./cmd/nox

# Build static binary (no CGO)
build-static:
	CGO_ENABLED=0 go build -o ./bin/$(BINARY) ./cmd/nox

# Cross-compile for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/$(BINARY)-linux-amd64 ./cmd/nox
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/$(BINARY)-darwin-amd64 ./cmd/nox
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/$(BINARY)-darwin-arm64 ./cmd/nox
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/$(BINARY)-windows-amd64.exe ./cmd/nox

# Clean build artifacts
clean:
	rm -f ./bin

# Run tests
test:
	go test -v ./...

# Run with test model
run: build
	./bin/$(BINARY) test assets/obj/nox.obj

# Install to GOPATH/bin
install:
	go install ./cmd/nox

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...
