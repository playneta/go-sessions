.PHONY:test
test:
	go vet ./... && golangci-lint run --tests=false && go test -race -v ./...

.PHONY: gen_mocks
gen_mocks:
	mockgen -source=./src/repositories/user.go -destination=./src/repositories/mocks/user.go
	mockgen -source=./src/providers/hash.go -destination=./src/providers/mocks/hash.go
	mockgen -source=./src/services/account.go -destination=./src/services/mocks/account.go

.PHONY: run
run:
	go run main.go