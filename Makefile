.PHONY: build down

build:
	docker-compose up
down:
	docker-compose down
	docker image prune -a -f
	rd /s /q db\postgresql\.database
stress:
	start "" "http://localhost:3003/d/8U8kBnaVz/dashboard?orgId=1&refresh=5s&from=now-10s&to=now"
	docker-compose run k6
