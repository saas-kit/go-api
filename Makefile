.PHONY: up down restart

up:
	cp .env.local .env
	docker-compose up -d --build
	docker logs saas_api

down:
	docker-compose down -v --rmi=local
	rm -vf .env

restart: down up

qa:
	# TODO: configure metalinter https://github.com/saas-kit/gometalinter
	gometalinter \
	    --vendor \
	    --deadline=60s \
			--cyclo-over=5 \
	    ./app/...
	go-cleanarch.v1 -ignore-tests ./app/...

test:
	go test -v ./...
