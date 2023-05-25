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