proto:
	protoc  --go_out=./ --go-grpc_out=./ \
              ./game.proto
