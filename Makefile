.PHONY: docker-up
docker-up:
	docker-compose -f docker/docker-compose.yml up

.PHONY: docker-down
docker-down:
	docker-compose -f docker/docker-compose.yml down
	docker system prune --volumes --force
.PHONY: up
up:
	go run -race cmd/elastic/main.go