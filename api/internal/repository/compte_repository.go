package repository

import (
	"bibliotheque-api/internal/models"

	"gorm.io/gorm"
)

type CompteRepository struct {
	db *gorm.DB
}

func NewCompteRepository(db *gorm.DB) *CompteRepository {
	return &CompteRepository{db: db}
}

func (r *CompteRepository) Create(compte *models.Compte) error {
	return r.db.Create(compte).Error
}

func (r *CompteRepository) Save(compte *models.Compte) error {
	return r.db.Save(compte).Error
}

func (r *CompteRepository) FindByEmail(email string) (*models.Compte, error) {
	var compte models.Compte
	err := r.db.Where("email = ?", email).First(&compte).Error
	return &compte, err
}

func (r *CompteRepository) FindByID(id uint) (*models.Compte, error) {
	var compte models.Compte
	err := r.db.Preload("Emprunts.Exemplaire").First(&compte, id).Error
	return &compte, err
}

// FindByIDLight charge le compte sans ses emprunts (pour auth/middleware).
func (r *CompteRepository) FindByIDLight(id uint) (*models.Compte, error) {
	var compte models.Compte
	err := r.db.First(&compte, id).Error
	return &compte, err
}

// FindAllAvecRetards retourne les comptes ayant au moins un emprunt en retard.
func (r *CompteRepository) FindAllAvecRetards() ([]models.Compte, error) {
	var comptes []models.Compte
	err := r.db.
		Joins("INNER JOIN emprunts ON emprunts.compte_id = comptes.id AND emprunts.rendu = false AND emprunts.date_limite < NOW() AND emprunts.deleted_at IS NULL").
		Preload("Emprunts", "rendu = false AND date_limite < NOW()").
		Group("comptes.id").
		Find(&comptes).Error
	return comptes, err
}
