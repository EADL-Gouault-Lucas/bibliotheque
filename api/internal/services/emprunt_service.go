package services

import (
	"fmt"
	"log"
	"time"

	"bibliotheque-api/internal/models"
	"bibliotheque-api/internal/repository"
)

type EmpruntService struct {
	empruntRepo    *repository.EmpruntRepository
	compteRepo     *repository.CompteRepository
	exemplaireRepo *repository.ExemplaireRepository
}

func NewEmpruntService(
	empruntRepo *repository.EmpruntRepository,
	compteRepo *repository.CompteRepository,
	exemplaireRepo *repository.ExemplaireRepository,
) *EmpruntService {
	return &EmpruntService{
		empruntRepo:    empruntRepo,
		compteRepo:     compteRepo,
		exemplaireRepo: exemplaireRepo,
	}
}

type CreateEmpruntInput struct {
	ExemplaireID uint `json:"exemplaire_id" binding:"required"`
}

// CreateEmprunt - R1, R2, R3, R4, R5, R9
func (s *EmpruntService) CreateEmprunt(compteID uint, input CreateEmpruntInput) (*models.Emprunt, error) {
	compte, err := s.compteRepo.FindByID(compteID)
	if err != nil {
		return nil, err
	}

	// R9 : les bibliothécaires ne peuvent pas emprunter
	if compte.IsBibliothecaire {
		return nil, ErrBibliothecaireEmprunt
	}

	// R1 : aucun emprunt en retard
	if compte.IsEnRetard() {
		return nil, ErrEmpruntEnRetard
	}

	exemplaire, err := s.exemplaireRepo.FindByID(input.ExemplaireID)
	if err != nil {
		return nil, err
	}

	// R4 : exemplaire disponible
	if exemplaire.EstEmprunte {
		return nil, ErrExemplaireIndisponible
	}

	// R2 : caution suffisante
	if compte.CautionRestante < exemplaire.Caution {
		return nil, ErrCautionInsuffisante
	}

	// R3 : pas de doublon sur la même ressource
	existants, err := s.empruntRepo.FindActiveByCompteAndLivre(compteID, exemplaire.LivreID)
	if err != nil {
		return nil, err
	}
	if len(existants) > 0 {
		return nil, ErrLivreDejaEmprunte
	}

	// R5 : durée de 15 jours
	now := time.Now()
	emprunt := &models.Emprunt{
		DateEmprunt:  now,
		DateLimite:   now.AddDate(0, 0, 15),
		CompteID:     compteID,
		ExemplaireID: exemplaire.ID,
	}

	exemplaire.EstEmprunte = true
	compte.CautionRestante -= exemplaire.Caution

	if err := s.empruntRepo.Create(emprunt); err != nil {
		return nil, err
	}
	if err := s.exemplaireRepo.Save(exemplaire); err != nil {
		return nil, err
	}
	if err := s.compteRepo.Save(compte); err != nil {
		return nil, err
	}

	return emprunt, nil
}

// RetournerExemplaire - R6 : remboursement intégral de la caution
func (s *EmpruntService) RetournerExemplaire(empruntID uint) (*models.Emprunt, error) {
	emprunt, err := s.empruntRepo.FindByID(empruntID)
	if err != nil {
		return nil, err
	}
	if emprunt.Rendu {
		return nil, ErrEmpruntDejaRendu
	}

	exemplaire, err := s.exemplaireRepo.FindByID(emprunt.ExemplaireID)
	if err != nil {
		return nil, err
	}

	compte, err := s.compteRepo.FindByIDLight(emprunt.CompteID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	emprunt.Rendu = true
	emprunt.DateRetour = &now
	exemplaire.EstEmprunte = false
	// R6 : remboursement intégral
	compte.CautionRestante += exemplaire.Caution

	if err := s.empruntRepo.Save(emprunt); err != nil {
		return nil, err
	}
	if err := s.exemplaireRepo.Save(exemplaire); err != nil {
		return nil, err
	}
	if err := s.compteRepo.Save(compte); err != nil {
		return nil, err
	}

	return emprunt, nil
}

// MesEmprunts retourne les emprunts d'un compte - R8
func (s *EmpruntService) MesEmprunts(compteID uint) ([]models.Emprunt, error) {
	return s.empruntRepo.FindByCompteID(compteID)
}

// ListRetards retourne tous les emprunts en retard - R10
func (s *EmpruntService) ListRetards() ([]models.Emprunt, error) {
	return s.empruntRepo.FindAllEnRetard()
}

// EnvoyerRappels - R21 : affiche dans la console les comptes en retard
func (s *EmpruntService) EnvoyerRappels() ([]models.Compte, error) {
	comptes, err := s.compteRepo.FindAllAvecRetards()
	if err != nil {
		return nil, err
	}
	for _, c := range comptes {
		log.Printf("[RAPPEL] %s\n",
			fmt.Sprintf("Compte %s (%s) — %d emprunt(s) en retard.", c.GetNomComplet(), c.Email, len(c.Emprunts)))
	}
	return comptes, nil
}
