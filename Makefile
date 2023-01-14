.PHONY: build
build:
	@echo "Building..."
	@mkdir -p build
	go build -o build/jordi ./cmd/jordi

.PHONY: test
test:
	@echo "Testing..."
	go test -v ./...

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf build

.PHONY: lint
lint:
	@echo "Linting..."
	golangci-lint run
