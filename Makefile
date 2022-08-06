build: go build -o entain-task cmd/server/main.go

lint: golangci-lint run
