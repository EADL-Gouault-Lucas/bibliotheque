package models

import "time"

type ICompte interface {
	IsEnRetard() bool
	GetNomComplet() string
}

type Compte struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	Email            string    `json:"email"`
	Prenom           string    `json:"prenom"`
	Nom              string    `json:"nom"`
	CautionRestante  float64   `json:"caution_restante"`
	IsBibliothecaire bool      `json:"is_bibliothecaire"`
	Emprunts         []Emprunt `json:"emprunts,omitempty"`
}

func (c *Compte) IsEnRetard() bool {
	for _, e := range c.Emprunts {
		if e.IsEnRetard() {
			return true
		}
	}
	return false
}

func (c *Compte) GetNomComplet() string {
	return c.Prenom + " " + c.Nom
}
