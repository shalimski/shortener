FROM golang:1.19.1 as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o shortener ./cmd/shortener/main.go

FROM alpine:3.16
EXPOSE 8080
COPY --from=builder /app/shortener /app/shortener

WORKDIR /app
CMD ["./shortener"]
