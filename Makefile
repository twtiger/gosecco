default: lint test

lint:
	golint ./...

test:
	go test -cover -v ./...

deps-dev:
	go get github.com/golang/lint/golint
	go get gopkg.in/check.v1

deps-dev-u:
	go get -u github.com/golang/lint/golint
	go get -u gopkg.in/check.v1
