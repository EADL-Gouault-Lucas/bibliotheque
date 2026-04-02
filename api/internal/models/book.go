package models

import "gorm.io/gorm"

type Livre struct {
	gorm.Model
	Titre       string       `json:"titre"      gorm:"not null"`
	CodeBarre   string       `json:"code_barre" gorm:"uniqueIndex"`
	CodeISBN    string       `json:"code_isbn"  gorm:"uniqueIndex"`
	Auteurs     []string     `json:"auteurs"    gorm:"serializer:json;not null"`
	Exemplaires []Exemplaire `json:"exemplaires,omitempty" gorm:"foreignKey:LivreID"`
}
