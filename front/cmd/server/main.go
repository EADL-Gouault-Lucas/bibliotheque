package main

import (
	"log"
	"os"

	"bibliotheque-front/internal/apiclient"
	"bibliotheque-front/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env.local") // priorité locale, ne fail pas si absent
	_ = godotenv.Load()             // valeurs par défaut depuis .env

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	api := apiclient.New(apiURL)

	authHandler := handlers.NewAuthHandler(api)
	catalogueHandler := handlers.NewCatalogueHandler(api)
	empruntHandler := handlers.NewEmpruntHandler(api)
	adminHandler := handlers.NewAdminHandler(api)

	r := gin.Default()
	r.Static("/static", "./web/static")

	// ── Public ────────────────────────────────────────────────────────────────
	r.GET("/", catalogueHandler.ShowCatalogue)
	r.GET("/login", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/register", authHandler.ShowRegister)
	r.POST("/register", authHandler.Register)
	r.GET("/logout", authHandler.Logout)

	// ── Utilisateur connecté ──────────────────────────────────────────────────
	auth := r.Group("/", handlers.RequireAuth())
	{
		auth.GET("/emprunts", empruntHandler.ShowMesEmprunts)
		auth.POST("/emprunts", empruntHandler.CreateEmprunt)
	}

	// ── Bibliothécaire ────────────────────────────────────────────────────────
	biblio := r.Group("/", handlers.RequireAuth(), handlers.RequireBibliothecaire())
	{
		biblio.GET("/admin/livres/nouveau", catalogueHandler.ShowNouveauLivre)
		biblio.POST("/admin/livres", catalogueHandler.CreateLivre)
		biblio.GET("/admin/livres/:id/exemplaires/nouveau", catalogueHandler.ShowNouvelExemplaire)
		biblio.POST("/admin/livres/:id/exemplaires", catalogueHandler.CreateExemplaire)

		biblio.GET("/admin/emprunts", adminHandler.ShowEmprunts)
		biblio.POST("/admin/emprunts/:id/retour", adminHandler.RetourExemplaire)
		biblio.POST("/admin/emprunts/rappels", adminHandler.EnvoyerRappels)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Frontend demarre sur :%s\n", port) // #nosec G706 -- port is a server-side env var
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Serveur arrêté : ", err)
	}
}
