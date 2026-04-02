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

// ListLivres - GET /api/v1/livres (public)
func (h *LivreHandler) ListLivres(c *gin.Context) {
	livres, err := h.livreSvc.ListLivres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, livres)
}

// CreateLivre - POST /api/v1/livres (bibliothécaire)
// R10, R17
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

// GetLivre - GET /api/v1/livres/:id (public)
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

// AddExemplaire - POST /api/v1/livres/:id/exemplaires (bibliothécaire)
// R10, R19, R20
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
