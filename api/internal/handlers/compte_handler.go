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

// Register godoc
// @Summary     Créer un compte
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body services.CreateCompteInput true "Informations du nouveau compte"
// @Success     201 {object} models.Compte
// @Failure     400 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/register [post]
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

// Login godoc
// @Summary     S'authentifier
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body services.LoginInput true "Identifiants"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/login [post]
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
