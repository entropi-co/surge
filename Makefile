dev-docker:
	@echo Starting development compose
	docker compose -f docker-compose.development.yml build surge-core --no-cache
	docker compose -f docker-compose.development.yml up -d
	@echo Shutdown development compose

dev-postgres-standalone:
	@echo Starting development PostgreSQL as standalone container
	docker compose -f docker-compose.postgres.yml up -d
	@echo Started in background

dev-postgres-standalone-stop:
	@echo Stopping standalone development PostgreSQL container
	docker compose -f docker-compose.postgres.yml down
	@echo Stopped

dev-postgres-standalone-reset:
	@echo Recreating database 'surge_development'
	@docker exec -t surge-postgres-standalone psql postgres://postgres:postgres@localhost/surge_development -c '\set AUTOCOMMIT on\drop database surge_development; create database surge_development;'
	@echo Recreated

dev-run:
	go build -o surge
	./surge

dev-logs:
	@echo Inspecting compose logs
	docker compose -f docker-compose.development.yml --ansi=always logs surge-core

dev-logsnocolor:
	@echo Inspecting compose logs without color
	docker compose -f docker-compose.development.yml logs