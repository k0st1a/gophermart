.DEFAULT_GOAL := build

.PHONY:build
build:
	go build -C ./cmd/gophermart/ -o gophermart

.PHONY:clean
clean:
	-rm -f ./cmd/gophermart/gophermart

.PHONY:statictest
statictest:
	go vet -vettool=$$(which statictest) ./...


test: build statictest
	go test -v ./...
