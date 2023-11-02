lint:
	@go fmt ./... && go vet ./...

.PHONY: test
test:
	@go test -v ./... -coverprofile cover.out

.PHONY: test-e2e
test-e2e:
	@go test -count=1 ./test/e2e/... --race

install:
	@go install -ldflags="-s -w" ./...