run:
	go build && ./douyin-backend-go
test:
	go test
lint:
	golangci-lint run