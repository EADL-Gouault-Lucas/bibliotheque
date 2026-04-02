package models

import "time"

type Livre struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	Titre       string       `json:"titre"`

	CodeISBN    string       `json:"code_isbn"`
	Auteurs     []string     `json:"auteurs"`
	Exemplaires []Exemplaire `json:"exemplaires,omitempty"`
}
