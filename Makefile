VERSION := 1.0

.PHONY: run
run:
	go run ./cmd/shortener/main.go

test:
	go test ./... -count=1 -cover

.PHONY: lint
lint:
	golangci-lint run ./... 

.PHONY: docker-build
docker-build:
	docker build . -t shalimski/shortener:${VERSION}