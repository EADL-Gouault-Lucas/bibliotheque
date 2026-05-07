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

// MesEmprunts godoc
// @Summary     Mes emprunts en cours
// @Tags        emprunts
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array}  models.Emprunt
// @Failure     401 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /emprunts [get]
func (h *EmpruntHandler) MesEmprunts(c *gin.Context) {
	compteID := c.MustGet("compte_id").(uint)
	emprunts, err := h.empruntSvc.MesEmprunts(compteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, emprunts)
}

// CreateEmprunt godoc
// @Summary     Emprunter un exemplaire
// @Tags        emprunts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body services.CreateEmpruntInput true "Exemplaire à emprunter"
// @Success     201 {object} models.Emprunt
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     422 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /emprunts [post]
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

// RetourExemplaire godoc
// @Summary     Retourner un exemplaire
// @Tags        emprunts
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "ID de l'emprunt"
// @Success     200 {object} models.Emprunt
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     403 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /emprunts/{id}/retour [put]
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

// ListActifs godoc
// @Summary     Lister tous les emprunts actifs
// @Tags        emprunts
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array}  models.Emprunt
// @Failure     401 {object} map[string]string
// @Failure     403 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /emprunts/actifs [get]
func (h *EmpruntHandler) ListActifs(c *gin.Context) {
	emprunts, err := h.empruntSvc.ListActifs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, emprunts)
}

// ListRetards godoc
// @Summary     Lister les emprunts en retard
// @Tags        emprunts
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array}  models.Emprunt
// @Failure     401 {object} map[string]string
// @Failure     403 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /emprunts/retards [get]
func (h *EmpruntHandler) ListRetards(c *gin.Context) {
	emprunts, err := h.empruntSvc.ListRetards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, emprunts)
}

// EnvoyerRappels godoc
// @Summary     Envoyer des rappels aux emprunteurs en retard
// @Tags        emprunts
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} map[string]interface{}
// @Failure     401 {object} map[string]string
// @Failure     403 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /emprunts/rappels [post]
func (h *EmpruntHandler) EnvoyerRappels(c *gin.Context) {
	comptes, err := h.empruntSvc.EnvoyerRappels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur interne"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rappels_envoyes": len(comptes), "comptes": comptes})
}
