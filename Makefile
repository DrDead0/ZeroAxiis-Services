run:
	docker compose up --build

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up --build

logs:
	docker compose logs -f

clean:
	docker compose down -v