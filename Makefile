proto:
	protoc  --go_out=./ --go-grpc_out=./ \
              ./game.proto

mock:
	mockgen -destination=tests/mocks/mock_stream.go -package=grpc google.golang.org/grpc BidiStreamingServer