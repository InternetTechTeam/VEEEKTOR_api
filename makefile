all:
	docker compose up --build

clean:
	docker system prune -a
