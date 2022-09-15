.PHONY: run
run:
	go run ./cmd/shortener/main.go

.PHONY: lint
lint:
	golangci-lint run ./... 
