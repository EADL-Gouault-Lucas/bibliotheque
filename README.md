# Bibliothèque — Projet CI/CD

Application de gestion de bibliothèque composée d'une API REST (Go/Gin) et d'un frontend web (Go/Gin + Templ).

---

## Démarrage rapide

### Prérequis

- [Go 1.25.8+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [templ](https://templ.guide/quick-start/installation)

### Mode développement local

```bash
# 1. Démarrer la base de données
make dev-db

# 2. Lancer l'API (dans un terminal)
make run-api

# 3. Lancer le frontend (dans un autre terminal)
make run-front
```

### Mode Docker complet

```bash
make docker-build   # Build les images
make docker-up      # Démarre tous les services
make docker-down    # Arrête tout
make docker-logs    # Affiche les logs
```

---

## URLs d'accès

| Service | Local | Docker |
|---|---|---|
| **Frontend** | http://localhost:3000 | http://localhost:3000 |
| **API** | http://localhost:8080 | http://localhost:8080 |
| **Base de données** | `localhost:5432` | `localhost:5432` |

### Endpoints API

| Méthode | URL | Description |
|---|---|---|
| `POST` | http://localhost:8080/api/v1/auth/register | Créer un compte |
| `POST` | http://localhost:8080/api/v1/auth/login | Se connecter |
| `GET` | http://localhost:8080/api/v1/livres | Liste des livres |
| `POST` | http://localhost:8080/api/v1/livres | Ajouter un livre *(bibliothécaire)* |
| `POST` | http://localhost:8080/api/v1/livres/:id/exemplaires | Ajouter un exemplaire *(bibliothécaire)* |
| `GET` | http://localhost:8080/api/v1/emprunts | Mes emprunts |
| `POST` | http://localhost:8080/api/v1/emprunts | Créer un emprunt |
| `PUT` | http://localhost:8080/api/v1/emprunts/:id/retour | Retourner un exemplaire *(bibliothécaire)* |
| `GET` | http://localhost:8080/api/v1/emprunts/retards | Liste des retards *(bibliothécaire)* |
| `POST` | http://localhost:8080/api/v1/emprunts/rappels | Envoyer les rappels *(bibliothécaire)* |

### Base de données

```
Host:     localhost
Port:     5432
User:     postgres
Password: postgres
Database: bibliotheque
```

---

## Commandes utiles

```bash
make help        # Liste toutes les commandes disponibles
make lint        # Lint (golangci-lint) sur API + frontend
make test        # Tests unitaires
make vuln        # Audit de sécurité (govulncheck)
make tidy        # go mod tidy sur les deux modules
make templ       # Régénère les fichiers depuis les templates .templ
make swag        # Régénère la documentation Swagger
```
