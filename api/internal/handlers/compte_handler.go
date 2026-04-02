package handlers

import (
	"errors"
	"net/http"

	"bibliotheque-api/internal/services"

	"github.com/gin-gonic/gin"
)

type CompteHandler struct {
	compteSvc *services.CompteService
	tokenSvc  *services.TokenService
}

func NewCompteHandler(compteSvc *services.CompteService, tokenSvc *services.TokenService) *CompteHandler {
	return &CompteHandler{compteSvc: compteSvc, tokenSvc: tokenSvc}
}

// Register - POST /api/v1/auth/register
// R11, R12 : email unique, tous les champs obligatoires
func (h *CompteHandler) Register(c *gin.Context) {
	var input services.CreateCompteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	compte, err := h.compteSvc.CreateCompte(input)
	if err != nil {
		if errors.Is(err, services.ErrEmailExistant) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusCreated, compte)
}

// Login - POST /api/v1/auth/login
// R16 : même message d'erreur quel que soit le problème (anti-énumération)
func (h *CompteHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	compte, err := h.compteSvc.Login(input)
	if err != nil {
		if errors.Is(err, services.ErrIdentifiantsInvalides) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	token := h.tokenSvc.Generate(compte.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "compte": compte})
}
