.PHONY: help \
        build build-api build-front \
        run-api run-front \
        test test-api test-front \
        test-api-unit test-api-db test-api-e2e \
        test-front-unit test-front-e2e \
        lint lint-api lint-front \
        tidy tidy-api tidy-front \
        templ \
        swag swag-fmt \
        db-sql db-clear seed-test-db \
        vuln vuln-api vuln-front \
        docker-up docker-down docker-build docker-logs \
        playwright-install playwright-test \
        clean

# ── Variables ─────────────────────────────────────────────────────────────────
GO        := go
API_DIR   := ./api
FRONT_DIR := ./front

# ── Help ──────────────────────────────────────────────────────────────────────
help: ## Affiche cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2}'

# ── Build ─────────────────────────────────────────────────────────────────────
build: build-api build-front ## Build les deux projets

build-api: ## Build l'API
	cd $(API_DIR) && $(GO) build ./...

build-front: templ ## Build le frontend (génère les templates templ d'abord)
	cd $(FRONT_DIR) && $(GO) build ./...

run-api: ## Lance l'API en local (charge .env.local en priorité)
	cd $(API_DIR) && $(GO) run ./cmd/server/main.go

run-front: templ ## Lance le frontend en local (charge .env.local en priorité)
	cd $(FRONT_DIR) && $(GO) run ./cmd/server/main.go

# ── Tests ─────────────────────────────────────────────────────────────────────
test: test-api test-front ## Exécute tous les tests

test-api: test-api-unit test-api-db test-api-e2e ## Tous les tests API

test-api-unit: ## Tests unitaires de l'API avec coverage
	cd $(API_DIR) && $(GO) test -tags unit -v -coverprofile=coverage.out ./...

test-api-db: ## Tests base de données de l'API
	cd $(API_DIR) && $(GO) test -tags db -v ./...

test-api-e2e: ## Tests e2e de l'API
	cd $(API_DIR) && $(GO) test -tags e2e -v ./...

test-front: test-front-unit test-front-e2e ## Tous les tests frontend

test-front-unit: ## Tests unitaires du frontend
	cd $(FRONT_DIR) && $(GO) test -tags unit -v ./...

test-front-e2e: ## Tests e2e du frontend (API doit être démarrée)
	cd $(FRONT_DIR) && $(GO) test -tags e2e -v ./tests/...

# ── Lint ──────────────────────────────────────────────────────────────────────
lint: lint-api lint-front ## Lint les deux projets

lint-api: ## Lint de l'API (golangci-lint)
	cd $(API_DIR) && golangci-lint run ./...

lint-front: ## Lint du frontend (golangci-lint)
	cd $(FRONT_DIR) && golangci-lint run ./...

# ── Dépendances ───────────────────────────────────────────────────────────────
tidy: tidy-api tidy-front ## go mod tidy sur les deux projets

tidy-api: ## go mod tidy — API
	cd $(API_DIR) && $(GO) mod tidy

tidy-front: ## go mod tidy — Frontend
	cd $(FRONT_DIR) && $(GO) mod tidy

# ── Génération ────────────────────────────────────────────────────────────────
templ: ## Génère les fichiers Go depuis les templates .templ
	cd $(FRONT_DIR) && templ generate

swag: ## Génère la documentation Swagger de l'API (docs/)
	cd $(API_DIR) && swag init -g cmd/server/main.go -o docs/ --parseInternal

swag-fmt: ## Formate les annotations Swagger
	cd $(API_DIR) && swag fmt

# ── Base de données ───────────────────────────────────────────────────────────
db-sql: ## Injecte le jeu de données de développement en base (Docker)
	docker exec -i bibliotheque-db psql -U postgres -d bibliotheque < $(API_DIR)/cmd/sql/seed.sql

db-clear: ## Vide toutes les tables (remet les séquences à zéro)
	docker exec -i bibliotheque-db psql -U postgres -d bibliotheque < $(API_DIR)/cmd/sql/truncate.sql

seed-test-db: ## Injecte le jeu de données dans la base de test CI (attend DB_HOST, DB_USER, DB_NAME)
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f $(API_DIR)/cmd/sql/seed.sql

# ── Audit sécurité ────────────────────────────────────────────────────────────
vuln: vuln-api vuln-front ## Audit des dépendances (govulncheck)

vuln-api: ## govulncheck sur l'API
	cd $(API_DIR) && govulncheck ./...

vuln-front: ## govulncheck sur le frontend
	cd $(FRONT_DIR) && govulncheck ./...

# ── Docker ────────────────────────────────────────────────────────────────────
docker-build: ## Build les images Docker
	docker compose build

docker-up: ## Démarre tous les services Docker en arrière-plan
	docker compose up -d

docker-down: ## Arrête et supprime les conteneurs
	docker compose down

docker-logs: ## Affiche les logs en temps réel
	docker compose logs -f

# ── Playwright ────────────────────────────────────────────────────────────────
playwright-install: ## Installe Playwright et les navigateurs
	npm install
	npx playwright install --with-deps

playwright-test: ## Exécute les tests Playwright (serveurs doivent être démarrés)
	npx playwright test

# ── Nettoyage ─────────────────────────────────────────────────────────────────
clean: ## Supprime les artefacts générés
	cd $(API_DIR) && rm -f coverage.out
	cd $(API_DIR) && rm -rf docs/
	cd $(FRONT_DIR) && find . -name '*_templ.go' -delete
