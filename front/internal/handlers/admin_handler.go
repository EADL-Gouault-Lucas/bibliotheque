package handlers

import (
	"fmt"
	"net/http"

	"bibliotheque-front/internal/apiclient"
	"bibliotheque-front/internal/session"
	"bibliotheque-front/web/templates"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	api *apiclient.Client
}

func NewAdminHandler(api *apiclient.Client) *AdminHandler {
	return &AdminHandler{api: api}
}

// ShowRetards - GET /admin/retards
func (h *AdminHandler) ShowRetards(c *gin.Context) {
	user, _ := session.Get(c)
	emprunts, err := h.api.ListRetards(user.Token)
	if err != nil {
		emprunts = nil
	}

	successMsg := c.Query("success")
	rappels := 0
	if q := c.Query("rappels"); q != "" {
		fmt.Sscanf(q, "%d", &rappels)
	}

	comp := templates.Retards(emprunts, user, successMsg, rappels)
	comp.Render(c.Request.Context(), c.Writer)
}

// RetourExemplaire - POST /admin/emprunts/:id/retour
func (h *AdminHandler) RetourExemplaire(c *gin.Context) {
	user, _ := session.Get(c)
	empruntID := parseUint(c.Param("id"))

	_, err := h.api.RetourExemplaire(user.Token, empruntID)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/retards?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/admin/retards?success=Exemplaire+marque+comme+rendu")
}

// EnvoyerRappels - POST /admin/rappels
func (h *AdminHandler) EnvoyerRappels(c *gin.Context) {
	user, _ := session.Get(c)

	resp, err := h.api.EnvoyerRappels(user.Token)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/retards?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/retards?rappels=%d", resp.RappelsEnvoyes))
}
