# Makefile для управления сервисом и БД

.PHONY: up down logs rm

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

rm:
	docker-compose down -v 