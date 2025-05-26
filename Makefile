.PHONY: test

test:
	go test -race -covermode=atomic -coverprofile=coverage.out $(go list ./... | grep -v vendor/)