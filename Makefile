GOPATH := $(go env GOPATH) 

default: test-acceptance

test-unit:
	go test -v -race -count=1 ./internal/...

test-acceptance:
	go test -v --race --count=1 ./acceptance_test/...

lint:
	golangci-lint run
	revive -config ./revive.toml
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

tools: tool-golangci-lint tool-fumpt tool-moq

tool-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -c bash -s -- -b ${GOPATH}/bin v1.50.1

tool-revive:
	go install github.com/mgechev/revive@master

tool-fumpt:
	go install mvdan.cc/gofumpt

tool-moq:
	go install github.com/matryer/moq

run:
	cd cmd && go run main.go

build:
	cd cmd && go build main.go

docker-build:
	scripts/compile_docker.sh

docker-run:
	docker run -t --name=coding-challenge -p 8080:80 -d coding-challenge

docker-logs:
	docker logs -f coding-challenge