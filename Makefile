#!make

up:
	docker-compose up &

stop:
	docker-compose stop

restart:
	docker-compose restart

logs:
	docker logs -f bank-transaction-app