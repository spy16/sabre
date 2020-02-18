VERSION="`git describe --abbrev=0 --tags`"
COMMIT="`git rev-list -1 --abbrev-commit HEAD`"

all: clean fmt test build

fmt:
	@echo "Formatting..."
	@goimports -l -w ./

install:
	@echo "Installing slang to GOBIN..."
	@go install -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT}" ./cmd/slang/

clean:
	@echo "Cleaning up..."
	@rm -rf ./bin
	@go mod tidy -v

test:
	@echo "Running tests..."
	@go test -cover ./...

test-verbose:
	@echo "Running tests..."
	@go test -v -cover ./...

benchmark:
	@echo "Running benchmarks..."
	@go test -benchmem -run="none" -bench="Benchmark.*" -v ./...

build:
	@echo "Building..."
	@mkdir -p ./bin
	@go build -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT}" -o ./bin/slang ./cmd/slang/
