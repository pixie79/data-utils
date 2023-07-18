coverage:
	go test -coverpkg=./... ./...

test:
	go test -v ./... ./...

tidy:
	goimports -v -w .


tools:
	go mod tidy
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/tools/cmd/godoc@latest

lint:
	golangci-lint -v run

docs:
	godoc -http=localhost:8080 -index -index_interval=10s -play -notes=BUG
