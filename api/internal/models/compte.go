package models

import "gorm.io/gorm"

type Compte struct {
	gorm.Model
	Email            string    `json:"email"             gorm:"uniqueIndex;not null"`
	Prenom           string    `json:"prenom"            gorm:"not null"`
	Nom              string    `json:"nom"               gorm:"not null"`
	MotDePasse       string    `json:"-"                 gorm:"not null"`
	CautionRestante  float64   `json:"caution_restante"`
	IsBibliothecaire bool      `json:"is_bibliothecaire" gorm:"default:false"`
	Emprunts         []Emprunt `json:"emprunts,omitempty" gorm:"foreignKey:CompteID"`
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
