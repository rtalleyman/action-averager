.PHONY: all clean

all: build run clean

build:
	go build -o averager-run main.go

run: build
	./averager-run

clean:
	rm averager-run