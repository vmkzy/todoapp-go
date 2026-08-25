include .env
export

export PROJECT_ROOT := $(CURDIR)
env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres
	
env-cleanup:
	@powershell -Command "$$ans = Read-Host 'Очистить все volume файлы окружения? Опасность утери данных. [y/N]'; \
	if ($$ans -match '^[yY]$$') { \
		docker compose down todoapp-postgres port-forwarder; \
		Remove-Item -Recurse -Force '$(PROJECT_ROOT)\out\pgdata' -ErrorAction SilentlyContinue; \
		Write-Host 'Файлы окружения очищены' \
	} else { \
		Write-Host 'Очистка окружения отменена' \
	}"

migrate-create:
	$(if $(strip $(seq)),,$(error Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init))
	@docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@$(MAKE) migrate-action action=up

migrate-down:
	@$(MAKE) migrate-action action=down

migrate-action:
	$(if $(strip $(action)),,$(error Отсутствует необходимый параметр action. Пример: make migrate-action action=up))
	@docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
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

logs-cleanup:
	@powershell -Command "$$ans = Read-Host 'Очистить все log файлы? Опасность утери логов. [y/N]'; \
	if ($$ans -eq 'y') { \
		Remove-Item -Recurse -Force '$(PROJECT_ROOT)\out\logs' -ErrorAction SilentlyContinue; \
		Write-Host 'Файлы логов очищены' \
	} else { \
		Write-Host 'Очистка логов отменена' \
	}"

todoapp-deploy:
	@docker compose up -d --build todoapp

todoapp-undeploy:
	@docker compose down todoapp
ps:
	@docker compose ps