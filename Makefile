lint:
	golangci-lint run -c ./golangci.yml

test:
	go test -cover ./... -count=1

test-bv:
	go test -cover ./bodyvalidator/... -count=1

test-hc: 
	go test -cover ./healthcheck/... -count=1

test-gs:
	go test -cover ./gracefulshutdown/... -count=1

coverage:
	go test -coverprofile coverage.out ./... && go tool cover -html=coverage.out

format:
	golines -w .
	go fmt ./...
