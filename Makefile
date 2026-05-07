.PHONY: help \
        build build-api build-front \
        run-api run-front \
        test test-api test-front \
        lint lint-api lint-front \
        tidy tidy-api tidy-front \
        templ \
        swag swag-fmt \
        vuln vuln-api \
        docker-up docker-down docker-build docker-logs \
        clean

# ── Variables ─────────────────────────────────────────────────────────────────
GO      := go
API_DIR := ./api
FRONT_DIR := ./front

# ── Help ──────────────────────────────────────────────────────────────────────
help: ## Affiche cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

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

test-api: ## Tests de l'API avec coverage
	cd $(API_DIR) && $(GO) test ./... -v -coverprofile=coverage.out

test-front: ## Tests du frontend
	cd $(FRONT_DIR) && $(GO) test ./... -v

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

# ── Audit sécurité ────────────────────────────────────────────────────────────
vuln: vuln-api ## Audit des dépendances (govulncheck)

vuln-api: ## govulncheck sur l'API
	cd $(API_DIR) && govulncheck ./...

# ── Docker ────────────────────────────────────────────────────────────────────
docker-build: ## Build les images Docker
	docker compose build

docker-up: ## Démarre tous les services Docker en arrière-plan
	docker compose up -d

docker-down: ## Arrête et supprime les conteneurs
	docker compose down

docker-logs: ## Affiche les logs en temps réel
	docker compose logs -f

# ── Nettoyage ─────────────────────────────────────────────────────────────────
clean: ## Supprime les artefacts générés
	cd $(API_DIR) && rm -f coverage.out
	cd $(API_DIR) && rm -rf docs/
	cd $(FRONT_DIR) && find . -name '*_templ.go' -delete
