lint:
	golangci-lint run -c ./golangci.yml

test:
	go test -cover ./... -count=1 -json | tparse --all

test-bv:
	go test -cover ./bodyvalidator/... -count=1 -json | tparse --all

test-hc: 
	go test -cover ./healthcheck/... -count=1 -json | tparse --all

test-gs:
	go test -cover ./gracefulshutdown/... -count=1 -json | tparse --all

coverage:
	go test -coverprofile coverage.out ./... && go tool cover -html=coverage.out

format:
	golines -w .
	go fmt ./...
