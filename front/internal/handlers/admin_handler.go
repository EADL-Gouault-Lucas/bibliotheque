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

// ShowEmprunts - GET /admin/emprunts
func (h *AdminHandler) ShowEmprunts(c *gin.Context) {
	user, _ := session.Get(c)
	emprunts, err := h.api.ListActifs(user.Token)
	if err != nil {
		emprunts = nil
	}

	hasRetards := false
	for i := range emprunts {
		if emprunts[i].IsEnRetard() {
			hasRetards = true
			break
		}
	}

	errMsg := c.Query("error")
	successMsg := c.Query("success")
	rappels := 0
	if q := c.Query("rappels"); q != "" {
		_, _ = fmt.Sscanf(q, "%d", &rappels)
	}

	comp := templates.Retards(emprunts, hasRetards, user, errMsg, successMsg, rappels)
	_ = comp.Render(c.Request.Context(), c.Writer)
}

// RetourExemplaire - POST /admin/emprunts/:id/retour
func (h *AdminHandler) RetourExemplaire(c *gin.Context) {
	user, _ := session.Get(c)
	empruntID := parseUint(c.Param("id"))

	_, err := h.api.RetourExemplaire(user.Token, empruntID)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/emprunts?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/admin/emprunts?success=Exemplaire+marque+comme+rendu")
}

// EnvoyerRappels - POST /admin/emprunts/rappels
func (h *AdminHandler) EnvoyerRappels(c *gin.Context) {
	user, _ := session.Get(c)

	resp, err := h.api.EnvoyerRappels(user.Token)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/emprunts?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/emprunts?rappels=%d", resp.RappelsEnvoyes))
}
