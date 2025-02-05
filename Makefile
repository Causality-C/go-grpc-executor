service:
	protoc --go_out=. --go-grpc_out=. --proto_path=protobufs protobufs/executor.proto

clean:
	rm -rf gen/