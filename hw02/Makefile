up:
	@docker-compose up -d

run: up
	@docker exec -it clickhouse-server sh

logs:
	@docker-compose logs clickhouse-server

down:
	@docker-compose down -v 
