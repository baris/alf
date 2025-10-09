NAME := alf
BIN  := bin/$(NAME)

.PHONY: build run test fmt tidy clean

build: $(BIN)

$(BIN):
	@echo "Building $(NAME)..."
	@mkdir -p bin
	@go build -o $(BIN) ./cmd/alf

run: build
	@$(BIN) $(ARGS)

test:
	@go test ./...

fmt:
	@gofmt -w $(shell go list -f '{{.Dir}}' ./...)

tidy:
	@go mod tidy

clean:
	@rm -rf bin
