.PHONY: up copy-env build run all
CONTAINER_API_CLIMA=api-clima-eco

up:
	docker compose up -d

copy-env:
	cp ./api-clima/.env.example ./api-clima/.env

run:
	docker exec $(CONTAINER_API_CLIMA) go run main.go

wait-db:
	@echo "Esperando banco de dados..."
	sleep 10
	@echo "Banco de dados pronto!"

start: up copy-env wait-db run