build:
	cd logs-handler && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/logs-handler

.PHONY: test
test:
	go test ./... -race -coverprofile coverage.out -covermode atomic

.PHONY: clean
clean:
	rm -rf ./bin
