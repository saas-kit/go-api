.PHONY: up down restart

up:
	cp .env.local .env
	docker-compose up -d --build
	docker logs saas_api

down:
	docker-compose down -v --rmi=local
	rm -vf .env

restart: down up
