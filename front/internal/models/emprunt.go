package models

import "time"

type IEmprunt interface {
	IsEnRetard() bool
}

type Emprunt struct {
	ID           uint       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	DateEmprunt  time.Time  `json:"date_emprunt"`
	DateLimite   time.Time  `json:"date_limite"`
	DateRetour   *time.Time `json:"date_retour"`
	Rendu        bool       `json:"rendu"`
	CompteID     uint       `json:"compte_id"`
	Compte       *Compte    `json:"compte,omitempty"`
	ExemplaireID uint       `json:"exemplaire_id"`
	Exemplaire   Exemplaire `json:"exemplaire,omitempty"`
}

func (e *Emprunt) IsEnRetard() bool {
	if e.Rendu {
		return false
	}
	return time.Now().After(e.DateLimite)
}
