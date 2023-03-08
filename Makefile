.PHONY: build down

build:
	docker-compose up
down:
	docker-compose down
	docker image prune -a -f
	rd /s /q db\postgresql\.database
stress:
	start "" "http://www.localhost:3003/d/oKgDrWa4k/dashboard?orgId=1&refresh=10s"
	docker-compose run k6

xz2134124:
	start "https://www.google.com"