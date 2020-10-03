#!make

DOCKER_COMPOSE_EXEC:= docker exec bank-transaction-app bash -c

up:
	docker-compose up &

stop:
	docker-compose stop

restart:
	docker-compose restart

logs:
	docker logs -f bank-transaction-app

test:
	$(DOCKER_COMPOSE_EXEC) 'go test -race -cover ./...'

# Example1: make test-cov-html PACKAGE=domain
# Example2: make test-cov-html PACKAGE=usecase
test-cov-html:
	$(DOCKER_COMPOSE_EXEC) 'go test -coverprofile cover.out ./$(PACKAGE) && \
	go tool cover -html=cover.out -o cover.html' && \
	xdg-open ./cover.html

clear:
	- sudo rm -rf ./.cover ./report ./main
	- sudo find . -name "*.html" -type f -delete
	- sudo find . -name "*.out" -type f -delete