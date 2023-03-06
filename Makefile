.PHONY: build down

build:
	docker-compose up
down:
	docker-compose down
	rd /s /q db\postgresql\.database