package repository

import (
	"bibliotheque-api/internal/models"

	"gorm.io/gorm"
)

type ExemplaireRepository struct {
	db *gorm.DB
}

func NewExemplaireRepository(db *gorm.DB) *ExemplaireRepository {
	return &ExemplaireRepository{db: db}
}

func (r *ExemplaireRepository) Create(exemplaire *models.Exemplaire) error {
	return r.db.Create(exemplaire).Error
}

func (r *ExemplaireRepository) FindByID(id uint) (*models.Exemplaire, error) {
	var exemplaire models.Exemplaire
	err := r.db.First(&exemplaire, id).Error
	return &exemplaire, err
}

func (r *ExemplaireRepository) Save(exemplaire *models.Exemplaire) error {
	return r.db.Save(exemplaire).Error
}
