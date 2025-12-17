.PHONY: dev prod down clean

dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d

down:
	docker compose down

# Clean everything
clean:
	docker compose down -v --rmi all
