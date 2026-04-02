package services

import "errors"

// Erreurs métier
var (
	// Emprunt
	ErrEmpruntEnRetard        = errors.New("le compte a des emprunts en retard")
	ErrCautionInsuffisante    = errors.New("caution restante insuffisante pour cet exemplaire")
	ErrExemplaireIndisponible = errors.New("cet exemplaire est déjà emprunté")
	ErrLivreDejaEmprunte      = errors.New("vous avez déjà un exemplaire de ce livre en cours d'emprunt")
	ErrBibliothecaireEmprunt  = errors.New("les bibliothécaires ne peuvent pas emprunter")
	ErrEmpruntDejaRendu       = errors.New("cet emprunt a déjà été retourné")

	// Compte
	ErrEmailExistant         = errors.New("cette adresse email est déjà utilisée")
	ErrIdentifiantsInvalides = errors.New("email ou mot de passe incorrect")

	// Livre / Exemplaire
	ErrLivreSansAuteur = errors.New("le livre doit avoir au moins un auteur")

	ErrLivreInexistant = errors.New("livre introuvable dans le catalogue")
)
