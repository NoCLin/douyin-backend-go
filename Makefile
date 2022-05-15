run:
	go build && ./douyin-backend-go
test:
	go test -v ./...
lint:
	golangci-lint run