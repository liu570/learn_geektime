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



mockgen:
	mockgen -destination=micro/rpc/mock_proxy_gen_test.go -package=rpc -source=micro/rpc/types.go Proxy
