package repository

import (
	"bibliotheque-api/internal/models"

	"gorm.io/gorm"
)

type LivreRepository struct {
	db *gorm.DB
}

func NewLivreRepository(db *gorm.DB) *LivreRepository {
	return &LivreRepository{db: db}
}

func (r *LivreRepository) Create(livre *models.Livre) error {
	return r.db.Create(livre).Error
}

func (r *LivreRepository) FindAll() ([]models.Livre, error) {
	var livres []models.Livre
	err := r.db.Preload("Exemplaires").Find(&livres).Error
	return livres, err
}

func (r *LivreRepository) FindByID(id uint) (*models.Livre, error) {
	var livre models.Livre
	err := r.db.Preload("Exemplaires").First(&livre, id).Error
	return &livre, err
}
