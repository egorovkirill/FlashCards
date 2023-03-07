.PHONY: build down

build:
	docker-compose up
down:
	docker-compose down
	docker image prune -a -f
	rd /s /q db\postgresql\.database
monitor:
	start http://localhost:3000