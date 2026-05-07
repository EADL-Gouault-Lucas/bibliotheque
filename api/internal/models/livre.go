package models

import (
	"time"

	"gorm.io/gorm"
)

type Livre struct {
	ID          uint           `json:"id"         gorm:"primarykey;autoIncrement"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"          gorm:"index" swaggerignore:"true"`
	Titre       string         `json:"titre"      gorm:"not null"`
	CodeISBN    string         `json:"code_isbn"  gorm:"uniqueIndex"`
	Auteurs     []string       `json:"auteurs"    gorm:"serializer:json;not null"`
	Exemplaires []Exemplaire   `json:"exemplaires,omitempty" gorm:"foreignKey:LivreID"`
}
