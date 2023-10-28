VERSION = v0.1.0

lint:
	@go fmt ./... && go vet ./...

.PHONY: test
test:
	@go test -v ./... -coverprofile cover.out

.PHONY: test-e2e
test-e2e:
	@go test -v -count=1 ./test/e2e --race

build: clean
	@go build  -ldflags="-s -w -X main.version=$(VERSION)" -o bin/kubectl-curl .

install:
	@go install -ldflags="-s -w -X main.version=$(VERSION)"

clean:
	@ rm -rf bin
