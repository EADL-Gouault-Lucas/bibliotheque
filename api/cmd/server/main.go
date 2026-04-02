package main

import (
	"log"
	"os"

	"bibliotheque-api/internal/handlers"
	"bibliotheque-api/internal/models"
	"bibliotheque-api/internal/repository"
	"bibliotheque-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Bibliothèque API
// @version         1.0
// @description     API REST pour la gestion d'une bibliothèque
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	godotenv.Load(".env.local") // priorité locale, ne fail pas si absent
	godotenv.Load()             // valeurs par défaut depuis .env

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Connexion BDD impossible : ", err)
	}

	// AutoMigrate : crée/met à jour les tables automatiquement
	db.AutoMigrate(
		&models.Compte{},
		&models.Livre{},
		&models.Exemplaire{},
		&models.Emprunt{},
	)

	// ── Repositories ──────────────────────────────────────────────────────────
	compteRepo := repository.NewCompteRepository(db)
	livreRepo := repository.NewLivreRepository(db)
	exemplaireRepo := repository.NewExemplaireRepository(db)
	empruntRepo := repository.NewEmpruntRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		tokenSecret = "dev-secret-change-in-production"
	}
	tokenSvc := services.NewTokenService(tokenSecret)
	compteSvc := services.NewCompteService(compteRepo)
	livreSvc := services.NewLivreService(livreRepo, exemplaireRepo)
	empruntSvc := services.NewEmpruntService(empruntRepo, compteRepo, exemplaireRepo)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authMiddleware := handlers.NewAuthMiddleware(tokenSvc, compteRepo)
	compteHandler := handlers.NewCompteHandler(compteSvc, tokenSvc)
	livreHandler := handlers.NewLivreHandler(livreSvc)
	empruntHandler := handlers.NewEmpruntHandler(empruntSvc)

	r := gin.Default()
	v1 := r.Group("/api/v1")

	// Auth (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", compteHandler.Register)
		auth.POST("/login", compteHandler.Login)
	}

	// Catalogue (lecture publique, écriture bibliothécaire)
	livres := v1.Group("/livres")
	{
		livres.GET("", livreHandler.ListLivres)
		livres.POST("", authMiddleware.Require(), authMiddleware.RequireBibliothecaire(), livreHandler.CreateLivre)
		livres.POST("/:id/exemplaires", authMiddleware.Require(), authMiddleware.RequireBibliothecaire(), livreHandler.AddExemplaire)
	}

	// Emprunts
	emprunts := v1.Group("/emprunts", authMiddleware.Require())
	{
		emprunts.GET("", empruntHandler.MesEmprunts)
		emprunts.POST("", empruntHandler.CreateEmprunt)
		// Routes bibliothécaire
		emprunts.PUT("/:id/retour", authMiddleware.RequireBibliothecaire(), empruntHandler.RetourExemplaire)
		emprunts.GET("/retards", authMiddleware.RequireBibliothecaire(), empruntHandler.ListRetards)
		emprunts.POST("/rappels", authMiddleware.RequireBibliothecaire(), empruntHandler.EnvoyerRappels)
	}

	log.Println("API démarrée sur :8080")
	r.Run(":8080")
}
