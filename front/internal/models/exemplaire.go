package models

import "time"

type Exemplaire struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	CodeBarre   string    `json:"code_barre"`
	EstEmprunte bool      `json:"est_emprunte"`
	Caution     float64   `json:"caution"`
	Travee      string    `json:"travee"`
	Etagere     string    `json:"etagere"`
	Niveau      string    `json:"niveau"`
	LivreID     uint      `json:"livre_id"`
	Livre       *Livre    `json:"livre,omitempty"`
}
