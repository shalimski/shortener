VERSION := 1.0

.PHONY: run
run:
	go run ./cmd/shortener/main.go

.PHONY: test
test:
	go test ./... -count=1 -cover

.PHONY: test-integration
test-integration:
	go test -tags=integration ./tests -count=1 -cover -coverpkg=./...

.PHONY: lint
lint:
	golangci-lint run ./... 

.PHONY: gen
gen:
	mockgen -source ./internal/ports/ports.go -destination ./internal/ports/mock/ports_mock.go	

.PHONY: docker-build
docker-build:
	docker build . -t shalimski/shortener:${VERSION}