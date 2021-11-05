lint:
	golangci-lint run -c ./golangci.yml

test-hc: 
	go test -cover ./healthcheck/... -count=1

format: 
	go fmt ./...