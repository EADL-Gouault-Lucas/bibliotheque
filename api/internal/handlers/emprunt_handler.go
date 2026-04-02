package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"bibliotheque-api/internal/services"

	"github.com/gin-gonic/gin"
)

type EmpruntHandler struct {
	empruntSvc *services.EmpruntService
}

func NewEmpruntHandler(empruntSvc *services.EmpruntService) *EmpruntHandler {
	return &EmpruntHandler{empruntSvc: empruntSvc}
}

// MesEmprunts - GET /api/v1/emprunts (utilisateur connecté)
// R8
func (h *EmpruntHandler) MesEmprunts(c *gin.Context) {
	compteID := c.MustGet("compte_id").(uint)
	emprunts, err := h.empruntSvc.MesEmprunts(compteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, emprunts)
}

// CreateEmprunt - POST /api/v1/emprunts (utilisateur connecté, non bibliothécaire)
// R1, R2, R3, R4, R5, R8, R9
func (h *EmpruntHandler) CreateEmprunt(c *gin.Context) {
	compteID := c.MustGet("compte_id").(uint)

	var input services.CreateEmpruntInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emprunt, err := h.empruntSvc.CreateEmprunt(compteID, input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBibliothecaireEmprunt),
			errors.Is(err, services.ErrEmpruntEnRetard),
			errors.Is(err, services.ErrCautionInsuffisante),
			errors.Is(err, services.ErrExemplaireIndisponible),
			errors.Is(err, services.ErrLivreDejaEmprunte):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		}
		return
	}

	c.JSON(http.StatusCreated, emprunt)
}

// RetourExemplaire - PUT /api/v1/emprunts/:id/retour (bibliothécaire)
// R6, R10
func (h *EmpruntHandler) RetourExemplaire(c *gin.Context) {
	empruntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identifiant invalide"})
		return
	}

	emprunt, err := h.empruntSvc.RetournerExemplaire(uint(empruntID))
	if err != nil {
		if errors.Is(err, services.ErrEmpruntDejaRendu) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}

	c.JSON(http.StatusOK, emprunt)
}

// ListActifs - GET /api/v1/emprunts/actifs (bibliothécaire)
func (h *EmpruntHandler) ListActifs(c *gin.Context) {
	emprunts, err := h.empruntSvc.ListActifs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, emprunts)
}

// ListRetards - GET /api/v1/emprunts/retards (bibliothécaire)
// R7, R10
func (h *EmpruntHandler) ListRetards(c *gin.Context) {
	emprunts, err := h.empruntSvc.ListRetards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, emprunts)
}

// EnvoyerRappels - POST /api/v1/emprunts/rappels (bibliothécaire)
// R21
func (h *EmpruntHandler) EnvoyerRappels(c *gin.Context) {
	comptes, err := h.empruntSvc.EnvoyerRappels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rappels_envoyes": len(comptes), "comptes": comptes})
}
