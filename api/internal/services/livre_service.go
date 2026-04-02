package services

import (
	"bibliotheque-api/internal/models"
	"bibliotheque-api/internal/repository"
)

type LivreService struct {
	livreRepo      *repository.LivreRepository
	exemplaireRepo *repository.ExemplaireRepository
}

func NewLivreService(livreRepo *repository.LivreRepository, exemplaireRepo *repository.ExemplaireRepository) *LivreService {
	return &LivreService{livreRepo: livreRepo, exemplaireRepo: exemplaireRepo}
}

type CreateLivreInput struct {
	Titre     string   `json:"titre"      binding:"required"`
	CodeISBN  string   `json:"code_isbn"  binding:"required"`
	Auteurs   []string `json:"auteurs"    binding:"required"`
}

type AddExemplaireInput struct {
	CodeBarre string  `json:"code_barre" binding:"required"`
	Caution   float64 `json:"caution"    binding:"required,gte=0"`
	Travee    string  `json:"travee"     binding:"required"`
	Etagere   string  `json:"etagere"    binding:"required"`
	Niveau    string  `json:"niveau"     binding:"required"`
}

func (s *LivreService) GetLivre(id uint) (*models.Livre, error) {
	livre, err := s.livreRepo.FindByID(id)
	if err != nil {
		return nil, ErrLivreInexistant
	}
	return livre, nil
}

// CreateLivre - R17 : titre, ISBN, au moins un auteur, au moins un thème
func (s *LivreService) CreateLivre(input CreateLivreInput) (*models.Livre, error) {
	if len(input.Auteurs) == 0 {
		return nil, ErrLivreSansAuteur
	}
	livre := &models.Livre{
		Titre:    input.Titre,
		CodeISBN:  input.CodeISBN,
		Auteurs:   input.Auteurs,
	}
	return livre, s.livreRepo.Create(livre)
}

// AddExemplaire - R19, R20 : code barre unique, caution, emplacement ; livre doit exister
func (s *LivreService) AddExemplaire(livreID uint, input AddExemplaireInput) (*models.Exemplaire, error) {
	// R20 : vérification que la ressource existe
	if _, err := s.livreRepo.FindByID(livreID); err != nil {
		return nil, ErrLivreInexistant
	}
	exemplaire := &models.Exemplaire{
		CodeBarre: input.CodeBarre,
		Caution:   input.Caution,
		Travee:    input.Travee,
		Etagere:   input.Etagere,
		Niveau:    input.Niveau,
		LivreID:   livreID,
	}
	return exemplaire, s.exemplaireRepo.Create(exemplaire)
}

func (s *LivreService) ListLivres() ([]models.Livre, error) {
	return s.livreRepo.FindAll()
}
