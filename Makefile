run-db:
	docker-compose up

run-db-cli:
	docker exec -it influx_server influx