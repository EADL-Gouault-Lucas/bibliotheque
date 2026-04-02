package session

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

const cookieName = "bib_session"

type User struct {
	CompteID         uint   `json:"compte_id"`
	Prenom           string `json:"prenom"`
	Nom              string `json:"nom"`
	Token            string `json:"token"`
	IsBibliothecaire bool   `json:"is_bibliothecaire"`
}

func Set(c *gin.Context, user User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	encoded := base64.URLEncoding.EncodeToString(data)
	c.SetCookie(cookieName, encoded, 86400, "/", "", false, true)
	return nil
}

func Get(c *gin.Context) (*User, bool) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return nil, false
	}
	data, err := base64.URLEncoding.DecodeString(cookie)
	if err != nil {
		return nil, false
	}
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, false
	}
	return &user, true
}

func Clear(c *gin.Context) {
	c.SetCookie(cookieName, "", -1, "/", "", false, true)
}
