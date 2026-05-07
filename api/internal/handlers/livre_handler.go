package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"bibliotheque-api/internal/services"

	"github.com/gin-gonic/gin"
)

type LivreHandler struct {
	livreSvc *services.LivreService
}

func NewLivreHandler(livreSvc *services.LivreService) *LivreHandler {
	return &LivreHandler{livreSvc: livreSvc}
}

// ListLivres godoc
// @Summary     Lister tous les livres
// @Tags        livres
// @Produce     json
// @Success     200 {array}  models.Livre
// @Failure     500 {object} map[string]string
// @Router      /livres [get]
func (h *LivreHandler) ListLivres(c *gin.Context) {
	livres, err := h.livreSvc.ListLivres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, livres)
}

// CreateLivre godoc
// @Summary     Créer un livre
// @Tags        livres
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body services.CreateLivreInput true "Données du livre"
// @Success     201 {object} models.Livre
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     403 {object} map[string]string
// @Failure     422 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /livres [post]
func (h *LivreHandler) CreateLivre(c *gin.Context) {
	var input services.CreateLivreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	livre, err := h.livreSvc.CreateLivre(input)
	if err != nil {
		if errors.Is(err, services.ErrLivreSansAuteur) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusCreated, livre)
}

// GetLivre godoc
// @Summary     Obtenir un livre par son ID
// @Tags        livres
// @Produce     json
// @Param       id path int true "ID du livre"
// @Success     200 {object} models.Livre
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /livres/{id} [get]
func (h *LivreHandler) GetLivre(c *gin.Context) {
	livreID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identifiant de livre invalide"})
		return
	}

	livre, err := h.livreSvc.GetLivre(uint(livreID))
	if err != nil {
		if errors.Is(err, services.ErrLivreInexistant) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, livre)
}

// AddExemplaire godoc
// @Summary     Ajouter un exemplaire à un livre
// @Tags        livres
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path int true "ID du livre"
// @Param       body body services.AddExemplaireInput true "Données de l'exemplaire"
// @Success     201 {object} models.Exemplaire
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     403 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /livres/{id}/exemplaires [post]
func (h *LivreHandler) AddExemplaire(c *gin.Context) {
	livreID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identifiant de livre invalide"})
		return
	}

	var input services.AddExemplaireInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exemplaire, err := h.livreSvc.AddExemplaire(uint(livreID), input)
	if err != nil {
		if errors.Is(err, services.ErrLivreInexistant) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusCreated, exemplaire)
}
