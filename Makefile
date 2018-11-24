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
	    --exclude="composite literal uses unkeyed fields" \
	    --exclude="should have comment or be unexported" \
	    --exclude="Errors unhandled" \
			--cyclo-over=5 \
	    ./...
	go-cleanarch.v1 -ignore-tests

test:
	go test -v ./...
