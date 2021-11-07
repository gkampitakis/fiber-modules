lint:
	golangci-lint run -c ./golangci.yml

test:
	go test -cover ./... -count=1

test-hc: 
	go test -cover ./healthcheck/... -count=1

test-gs:
	go test -cover ./gracefulshutdown/... -count=1

format: 
	go fmt ./...