package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"bibliotheque-api/internal/models"
	"bibliotheque-api/internal/repository"

	"gorm.io/gorm"
)

type CompteService struct {
	repo *repository.CompteRepository
}

func NewCompteService(repo *repository.CompteRepository) *CompteService {
	return &CompteService{repo: repo}
}

type CreateCompteInput struct {
	Email      string `json:"email"        binding:"required,email"`
	Prenom     string `json:"prenom"       binding:"required"`
	Nom        string `json:"nom"          binding:"required"`
	MotDePasse string `json:"mot_de_passe" binding:"required,min=8"`
}

type LoginInput struct {
	Email      string `json:"email"       binding:"required,email"`
	MotDePasse string `json:"mot_de_passe" binding:"required"`
}

// CreateCompte - R11, R12
func (s *CompteService) CreateCompte(input CreateCompteInput) (*models.Compte, error) {
	// R11 : email unique
	existing, err := s.repo.FindByEmail(input.Email)
	if err == nil && existing.ID != 0 {
		return nil, ErrEmailExistant
	}

	hash, err := hashPassword(input.MotDePasse)
	if err != nil {
		return nil, err
	}

	compte := &models.Compte{
		Email:           input.Email,
		Prenom:          input.Prenom,
		Nom:             input.Nom,
		MotDePasse:      hash,
		CautionRestante: 0,
	}
	if err := s.repo.Create(compte); err != nil {
		return nil, err
	}
	return compte, nil
}

// Login - R16 (même message quel que soit l'erreur, anti-énumération)
func (s *CompteService) Login(input LoginInput) (*models.Compte, error) {
	compte, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdentifiantsInvalides
		}
		return nil, err
	}
	if !checkPassword(input.MotDePasse, compte.MotDePasse) {
		return nil, ErrIdentifiantsInvalides
	}
	return compte, nil
}

// hashPassword retourne sha256(salt+password) avec le sel en préfixe.
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	saltHex := hex.EncodeToString(salt)
	hash := sha256.Sum256([]byte(saltHex + password))
	return fmt.Sprintf("%s:%s", saltHex, hex.EncodeToString(hash[:])), nil
}

func checkPassword(password, stored string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return false
	}
	hash := sha256.Sum256([]byte(parts[0] + password))
	return hex.EncodeToString(hash[:]) == parts[1]
}
