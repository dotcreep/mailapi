.PHONY: docs
docs:
	@swag init -g ./cmd/main.go && swag fmt
	
.PHONY: run
run:
	@go run ./cmd/.

.PHONY: build
build:
	@go build -o ./bin/main ./cmd