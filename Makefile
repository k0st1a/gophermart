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

.PHONY:test
test: build statictest
	go test -v ./...

.PHONY:gophermarttest
gophermarttest: test
	./cmd/test/gophermarttest \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
		-accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
		-accrual-host=localhost \
		-accrual-port=$$(random unused-port) \
		-accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"

