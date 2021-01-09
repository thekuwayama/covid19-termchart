.PHONY: fmt
fmt:
	go fmt

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run covid19-termchart.go ${ARGS}

build:
	go build

.PHONY: install
install: fmt lint test
	go install

clean:
	go clean

.PHONY: all
all: install
