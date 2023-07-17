coverage:
	go test -coverpkg=./... ./...

test:
	go test -v ./... ./...

tidy:
	goimports -v -w .


tools:
	go mod tidy
	go install golang.org/x/tools/cmd/goimports@latest

lint:
	golangci-lint -v run
