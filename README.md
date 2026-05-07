# Bibliothèque — Projet CI/CD

Application de gestion de bibliothèque composée d'une API REST (Go/Gin) et d'un frontend web (Go/Gin + Templ).

---

## Démarrage rapide

### Prérequis

- [Go 1.25.9+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [templ](https://templ.guide/quick-start/installation)
- [swag](https://github.com/swaggo/swag) — `go install github.com/swaggo/swag/cmd/swag@latest`

### Lancer le projet complet

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
