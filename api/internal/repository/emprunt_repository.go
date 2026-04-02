package repository

import (
	"bibliotheque-api/internal/models"

	"gorm.io/gorm"
)

type EmpruntRepository struct {
	db *gorm.DB
}

func NewEmpruntRepository(db *gorm.DB) *EmpruntRepository {
	return &EmpruntRepository{db: db}
}

func (r *EmpruntRepository) Create(emprunt *models.Emprunt) error {
	return r.db.Create(emprunt).Error
}

func (r *EmpruntRepository) Save(emprunt *models.Emprunt) error {
	return r.db.Save(emprunt).Error
}

func (r *EmpruntRepository) FindByID(id uint) (*models.Emprunt, error) {
	var emprunt models.Emprunt
	err := r.db.Preload("Exemplaire.Livre").First(&emprunt, id).Error
	return &emprunt, err
}

func (r *EmpruntRepository) FindByCompteID(compteID uint) ([]models.Emprunt, error) {
	var emprunts []models.Emprunt
	err := r.db.Preload("Exemplaire.Livre").Where("compte_id = ?", compteID).Find(&emprunts).Error
	return emprunts, err
}

// FindActiveByCompteAndLivre retourne les emprunts actifs (non rendus) pour un compte et un livre donnés
// afin de vérifier la règle R3 (pas de doublon sur la même ressource).
func (r *EmpruntRepository) FindActiveByCompteAndLivre(compteID, livreID uint) ([]models.Emprunt, error) {
	var emprunts []models.Emprunt
	err := r.db.
		Joins("JOIN exemplaires ON exemplaires.id = emprunts.exemplaire_id AND exemplaires.deleted_at IS NULL").
		Where("emprunts.compte_id = ? AND exemplaires.livre_id = ? AND emprunts.rendu = false", compteID, livreID).
		Find(&emprunts).Error
	return emprunts, err
}

// FindAllActifs retourne tous les emprunts non rendus.
func (r *EmpruntRepository) FindAllActifs() ([]models.Emprunt, error) {
	var emprunts []models.Emprunt
	err := r.db.Preload("Exemplaire.Livre").Preload("Compte").
		Where("rendu = false").
		Find(&emprunts).Error
	return emprunts, err
}

// FindAllEnRetard retourne tous les emprunts en retard (non rendus, date_limite dépassée).
func (r *EmpruntRepository) FindAllEnRetard() ([]models.Emprunt, error) {
	var emprunts []models.Emprunt
	err := r.db.Preload("Exemplaire.Livre").Preload("Compte").
		Where("rendu = false AND date_limite < NOW()").
		Find(&emprunts).Error
	return emprunts, err
}
