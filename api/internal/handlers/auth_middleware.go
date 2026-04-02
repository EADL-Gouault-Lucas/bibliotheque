package handlers

import (
	"net/http"
	"strings"

	"bibliotheque-api/internal/repository"
	"bibliotheque-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenSvc   *services.TokenService
	compteRepo *repository.CompteRepository
}

func NewAuthMiddleware(tokenSvc *services.TokenService, compteRepo *repository.CompteRepository) *AuthMiddleware {
	return &AuthMiddleware{tokenSvc: tokenSvc, compteRepo: compteRepo}
}

// Require vérifie le token et injecte le compte dans le contexte Gin. - R8
func (m *AuthMiddleware) Require() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentification requise"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		compteID, err := m.tokenSvc.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide"})
			return
		}

		compte, err := m.compteRepo.FindByIDLight(compteID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "compte introuvable"})
			return
		}

		c.Set("compte_id", compte.ID)
		c.Set("is_bibliothecaire", compte.IsBibliothecaire)
		c.Next()
	}
}

// RequireBibliothecaire vérifie que l'utilisateur connecté est bibliothécaire. - R10
func (m *AuthMiddleware) RequireBibliothecaire() gin.HandlerFunc {
	return func(c *gin.Context) {
		isBiblio, _ := c.Get("is_bibliothecaire")
		if isBiblio != true {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "accès réservé aux bibliothécaires"})
			return
		}
		c.Next()
	}
}


