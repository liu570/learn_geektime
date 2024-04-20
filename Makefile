start_otel:
	@docker version


# e2e 测试

e2e:
	service docker stop
	service docker start
	-docker compose down
	docker compose up -d
	go test -race ./orm/... -tags=e2e
	docker compose down
	service docker stop



mock_gen:
	mockgen -destination=micro/rpc/mock_proxy_gen_test.go -package=rpc -source=micro/rpc/types.go Proxy

mock_proto:
	# 需要下载 proto 的编译器  apt install protobuf-compiler
	#go install google.golang.org/protobuf/cmd/protoc-gen-go
	#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	protoc --go_out=. user.proto
	# 生成 grpc 的代码，（与上面代码不能合成一行，不清楚原因）
	protoc --go-grpc_out=. user.proto



re_init_go_env:
	cd $GOPATH
	rm -rf pkg bin
	mkdir pkg bin
	go mod tidy
	go mod download


