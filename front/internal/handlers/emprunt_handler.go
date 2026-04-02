package handlers

import (
	"net/http"

	"bibliotheque-front/internal/apiclient"
	"bibliotheque-front/internal/session"
	"bibliotheque-front/web/templates"

	"github.com/gin-gonic/gin"
)

type EmpruntHandler struct {
	api *apiclient.Client
}

func NewEmpruntHandler(api *apiclient.Client) *EmpruntHandler {
	return &EmpruntHandler{api: api}
}

// ShowMesEmprunts - GET /emprunts
func (h *EmpruntHandler) ShowMesEmprunts(c *gin.Context) {
	user, _ := session.Get(c)
	emprunts, err := h.api.MesEmprunts(user.Token)
	if err != nil {
		emprunts = nil
	}
	comp := templates.MesEmprunts(emprunts, user, c.Query("success"))
	comp.Render(c.Request.Context(), c.Writer)
}

// CreateEmprunt - POST /emprunts
func (h *EmpruntHandler) CreateEmprunt(c *gin.Context) {
	user, _ := session.Get(c)
	exemplaireID := parseUint(c.PostForm("exemplaire_id"))

	_, err := h.api.CreateEmprunt(user.Token, apiclient.CreateEmpruntInput{
		ExemplaireID: exemplaireID,
	})
	if err != nil {
		c.Redirect(http.StatusFound, "/?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/emprunts?success=Emprunt+enregistre+avec+succes")
}
