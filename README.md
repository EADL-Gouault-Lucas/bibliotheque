# Bibliothèque — Projet CI/CD

Application de gestion de bibliothèque composée d'une API REST (Go/Gin) et d'un frontend web (Go/Gin + Templ).

---

## Démarrage rapide

### Prérequis

- [Go 1.25.11+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [templ](https://templ.guide/quick-start/installation) — `go install github.com/a-h/templ/cmd/templ@latest`
- [swag](https://github.com/swaggo/swag) — `go install github.com/swaggo/swag/cmd/swag@latest`

### Lancer le projet avec Docker

```bash
# 1. Démarrer tous les services (API, frontend, base de données)
make docker-up

# 2. Injecter les données de démonstration
make db-sql
```

Le projet est ensuite accessible sur :

| Service | URL |
|---|---|
| **Frontend** | http://localhost:3000 |
| **API** | http://localhost:8080 |
| **Swagger UI** | http://localhost:8080/swagger/index.html |
| **Base de données** | `localhost:5432` |

Pour arrêter :

```bash
make docker-down
```

### Lancer le projet en local (développement)

```bash
# 1. Générer les fichiers depuis les templates .templ
make templ

# 2. Générer la documentation Swagger
make swag

# 3. Démarrer la base de données
make docker-up

# 4. Lancer l'API (dans un terminal)
make run-api

# 5. Lancer le frontend (dans un autre terminal)
make run-front
```

---

## Endpoints API

La documentation complète est disponible via Swagger UI au démarrage.

### Base de données

```
Host:     localhost
Port:     5432
User:     postgres
Password: postgres
Database: bibliotheque
```

---

### Images Docker

Les images sont publiées sur GitHub Container Registry à chaque release.

### Artefacts produits

Les résultats du scan Trivy sont disponibles dans l'onglet **Security → Code scanning** du dépôt, et le tableau récapitulatif dans l'onglet **Summary** de chaque run de publication.
