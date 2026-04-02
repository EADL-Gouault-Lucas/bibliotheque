package services

import "errors"

// Erreurs métier
var (
	// Emprunt
	ErrEmpruntEnRetard        = errors.New("Le compte a des emprunts en retard")
	ErrCautionInsuffisante    = errors.New("Caution restante insuffisante pour cet exemplaire")
	ErrExemplaireIndisponible = errors.New("Cet exemplaire est déjà emprunté")
	ErrLivreDejaEmprunte      = errors.New("Vous avez déjà un exemplaire de ce livre en cours d'emprunt")
	ErrBibliothecaireEmprunt  = errors.New("Les bibliothécaires ne peuvent pas emprunter")
	ErrEmpruntDejaRendu       = errors.New("Cet emprunt a déjà été retourné")

	// Compte
	ErrEmailExistant         = errors.New("Cette adresse email est déjà utilisée")
	ErrIdentifiantsInvalides = errors.New("Email ou mot de passe incorrect")

	// Livre / Exemplaire
	ErrLivreSansAuteur = errors.New("Le livre doit avoir au moins un auteur")

	ErrLivreInexistant = errors.New("Livre introuvable dans le catalogue")
)
