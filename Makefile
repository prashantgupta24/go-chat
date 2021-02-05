mockgen:
	mockgen -package mock_service github.com/go-chat/server \
		Connection > mock_service/Connection.go

start:
	go run main.go

test:
	go test ./... -v -failfast