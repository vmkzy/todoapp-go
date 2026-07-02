include .env
export

export PROJECT_ROOT := $(CURDIR)
env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres
	
env-cleanup:
	@powershell -Command "$$ans = Read-Host 'Clean all volume files [y/N]'; if ($$ans -eq 'y') \
	{ docker compose down todoapp-postgres port-forwarder; Remove-Item -Recurse -Force out\pgdata -ErrorAction SilentlyContinue; Write-Host 'File environment clean' } \
	else { Write-Host 'Clean environment cancel' }"

migrate-create:
	@powershell -Command "if ('$(seq)' -eq '') { Write-Host 'Erorr seq. Example make migrate-create seq=init'; exit 1 }"
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"
	
migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@powershell -Command "if ('$(action)' -eq '') { Write-Host 'Erorr action. Example make migrate-action action=up'; exit 1 }"
	@docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

todoapp-run:
	@set "LOGGER_FOLDER=%PROJECT_ROOT%\out\logs" && \
	set "POSTGRES_HOST=localhost" &&\
	go mod tidy && \
	go run cmd/todoapp/main.go
