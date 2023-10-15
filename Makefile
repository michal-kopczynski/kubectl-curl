VERSION = v0.1.0

lint:
	@go fmt ./... && go vet ./...

test:
	@go test -v ./... -coverprofile cover.out

build: clean
	@go build  -ldflags="-s -w -X main.version=$(VERSION)" -o bin/kubectl-curl .

install:
	@go install -ldflags="-s -w -X main.version=$(VERSION)"

clean:
	@ rm -rf bin
