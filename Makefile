start:
	docker compose -f task/docker-compose.yaml up -d --build && \
	docker compose -f auth/docker-compose.yaml up -d --build && \
	docker compose -f docker-compose.yaml up -d
stop:
	docker compose -f task/docker-compose.yaml down -v && \
	docker compose -f auth/docker-compose.yaml down -v && \
	docker compose -f docker-compose.yaml down -v
create-net:	
	docker network create app