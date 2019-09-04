.PHONY: all clean

all: build run clean

init:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

build:
	go build -o averager-run main.go

run: build
	./averager-run

clean:
	rm averager-run