package handlers

import (
	"net/http"
	"strings"

	"bibliotheque-front/internal/apiclient"
	"bibliotheque-front/internal/session"
	"bibliotheque-front/web/templates"

	"github.com/gin-gonic/gin"
)

type CatalogueHandler struct {
	api *apiclient.Client
}

func NewCatalogueHandler(api *apiclient.Client) *CatalogueHandler {
	return &CatalogueHandler{api: api}
}

// ShowCatalogue - GET /
func (h *CatalogueHandler) ShowCatalogue(c *gin.Context) {
	livres, err := h.api.ListLivres()
	errMsg := ""
	if err != nil {
		errMsg = "Impossible de charger le catalogue : " + err.Error()
		livres = nil
	}
	if q := c.Query("error"); q != "" {
		errMsg = q
	}

	user, _ := session.Get(c)
	comp := templates.Catalogue(livres, user, errMsg)
	_ = comp.Render(c.Request.Context(), c.Writer)
}

// ShowNouveauLivre - GET /livres/nouveau (biblio only)
func (h *CatalogueHandler) ShowNouveauLivre(c *gin.Context) {
	user, _ := session.Get(c)
	comp := templates.NouveauLivre(user, c.Query("error"))
	_ = comp.Render(c.Request.Context(), c.Writer)
}

// CreateLivre - POST /livres (biblio only)
func (h *CatalogueHandler) CreateLivre(c *gin.Context) {
	user, _ := session.Get(c)

	auteurs := splitAndTrim(c.PostForm("auteurs"))

	input := apiclient.CreateLivreInput{
		Titre:     c.PostForm("titre"),
		CodeBarre: c.PostForm("code_barre"),
		CodeISBN:  c.PostForm("code_isbn"),
		Auteurs:   auteurs,
	}

	_, err := h.api.CreateLivre(user.Token, input)
	if err != nil {
		c.Redirect(http.StatusFound, "/livres/nouveau?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// ShowNouvelExemplaire - GET /livres/:id/exemplaires/nouveau (biblio only)
func (h *CatalogueHandler) ShowNouvelExemplaire(c *gin.Context) {
	livreID := parseUint(c.Param("id"))
	livre, err := h.api.GetLivre(livreID)
	user, _ := session.Get(c)

	errMsg := c.Query("error")
	if err != nil {
		errMsg = "Livre introuvable."
	}

	if livre == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	comp := templates.NouvelExemplaire(*livre, user, errMsg)
	_ = comp.Render(c.Request.Context(), c.Writer)
}

// CreateExemplaire - POST /livres/:id/exemplaires (biblio only)
func (h *CatalogueHandler) CreateExemplaire(c *gin.Context) {
	livreID := parseUint(c.Param("id"))
	user, _ := session.Get(c)

	caution := parseFloat(c.PostForm("caution"))

	input := apiclient.AddExemplaireInput{
		CodeBarre: c.PostForm("code_barre"),
		Caution:   caution,
		Travee:    c.PostForm("travee"),
		Etagere:   c.PostForm("etagere"),
		Niveau:    c.PostForm("niveau"),
	}

	_, err := h.api.AddExemplaire(user.Token, livreID, input)
	if err != nil {
		redirect := "/livres/" + c.Param("id") + "/exemplaires/nouveau?error=" + encodeMsg(err.Error())
		c.Redirect(http.StatusFound, redirect)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
