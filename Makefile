.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/mylisp

.PHONY: run
run: build
	bin/mylisp

.PHONY: test
test: build
	go test -count=1 ./...

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rf bin
