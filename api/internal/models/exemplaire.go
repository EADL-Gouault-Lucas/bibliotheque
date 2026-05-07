package models

import (
	"time"

	"gorm.io/gorm"
)

type Exemplaire struct {
	ID          uint           `json:"id"         gorm:"primarykey;autoIncrement"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"          gorm:"index" swaggerignore:"true"`
	CodeBarre   string         `json:"code_barre"   gorm:"uniqueIndex;not null"`
	EstEmprunte bool           `json:"est_emprunte" gorm:"default:false"`
	Caution     float64        `json:"caution"     gorm:"not null"`
	Travee      string         `json:"travee"       gorm:"not null"`
	Etagere     string         `json:"etagere"      gorm:"not null"`
	Niveau      string         `json:"niveau"       gorm:"not null"`
	LivreID     uint           `json:"livre_id"     gorm:"not null"`
	Livre       *Livre         `json:"livre,omitempty" gorm:"foreignKey:LivreID"`
}
