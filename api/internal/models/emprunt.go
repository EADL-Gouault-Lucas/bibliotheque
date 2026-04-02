package models

import (
	"time"

	"gorm.io/gorm"
)

type Emprunt struct {
	gorm.Model
	DateEmprunt  time.Time  `json:"date_emprunt"  gorm:"not null"`
	DateLimite   time.Time  `json:"date_limite"   gorm:"not null"`
	DateRetour   *time.Time `json:"date_retour"`
	Rendu        bool       `json:"rendu"         gorm:"default:false"`
	CompteID     uint       `json:"compte_id"     gorm:"not null"`
	Compte       *Compte    `json:"compte,omitempty" gorm:"foreignKey:CompteID"`
	ExemplaireID uint       `json:"exemplaire_id" gorm:"not null"`
	Exemplaire   Exemplaire `json:"exemplaire,omitempty"`
}

func (e *Emprunt) IsEnRetard() bool {
	if e.Rendu {
		return false
	}
	return time.Now().After(e.DateLimite)
}
