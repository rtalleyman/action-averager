.PHONY: all test test-debug clean

all: build run test clean

init:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

build:
	go build -o averager-run main.go

run: build
	./averager-run

test:
	ginkgo -p -keepGoing --randomizeAllSpecs pkg/test/

test-debug:
	ginkgo -p -v -race -keepGoing --randomizeAllSpecs --progress --trace pkg/test/

clean:
	rm averager-run