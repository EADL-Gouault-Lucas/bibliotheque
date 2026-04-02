package handlers

import (
	"net/http"

	"bibliotheque-front/internal/apiclient"
	"bibliotheque-front/internal/session"
	"bibliotheque-front/web/templates"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	api *apiclient.Client
}

func NewAuthHandler(api *apiclient.Client) *AuthHandler {
	return &AuthHandler{api: api}
}

// ShowLogin - GET /login
func (h *AuthHandler) ShowLogin(c *gin.Context) {
	errMsg := c.Query("error")
	successMsg := c.Query("success")
	comp := templates.Login(errMsg, successMsg)
	_ = comp.Render(c.Request.Context(), c.Writer)
}

// Login - POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	input := apiclient.LoginInput{
		Email:      c.PostForm("email"),
		MotDePasse: c.PostForm("mot_de_passe"),
	}

	resp, err := h.api.Login(input)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error="+encodeMsg(err.Error()))
		return
	}

	_ = session.Set(c, session.User{
		CompteID:         resp.Compte.ID,
		Prenom:           resp.Compte.Prenom,
		Nom:              resp.Compte.Nom,
		Token:            resp.Token,
		IsBibliothecaire: resp.Compte.IsBibliothecaire,
	})
	c.Redirect(http.StatusFound, "/")
}

// ShowRegister - GET /register
func (h *AuthHandler) ShowRegister(c *gin.Context) {
	errMsg := c.Query("error")
	comp := templates.Register(errMsg)
	_ = comp.Render(c.Request.Context(), c.Writer)
}

// Register - POST /register
func (h *AuthHandler) Register(c *gin.Context) {
	input := apiclient.RegisterInput{
		Email:      c.PostForm("email"),
		Prenom:     c.PostForm("prenom"),
		Nom:        c.PostForm("nom"),
		MotDePasse: c.PostForm("mot_de_passe"),
	}

	_, err := h.api.Register(input)
	if err != nil {
		c.Redirect(http.StatusFound, "/register?error="+encodeMsg(err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/login?success=Compte+cree+avec+succes")
}

// Logout - GET /logout
func (h *AuthHandler) Logout(c *gin.Context) {
	session.Clear(c)
	c.Redirect(http.StatusFound, "/")
}
