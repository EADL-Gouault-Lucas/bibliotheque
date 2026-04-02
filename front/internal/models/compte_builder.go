package models

type CompteBuilder struct {
	compte Compte
}

func NewCompteBuilder() *CompteBuilder {
	return &CompteBuilder{}
}

func (b *CompteBuilder) SetID(id uint) *CompteBuilder {
	b.compte.ID = id
	return b
}

func (b *CompteBuilder) SetEmail(email string) *CompteBuilder {
	b.compte.Email = email
	return b
}

func (b *CompteBuilder) SetPrenom(prenom string) *CompteBuilder {
	b.compte.Prenom = prenom
	return b
}

func (b *CompteBuilder) SetNom(nom string) *CompteBuilder {
	b.compte.Nom = nom
	return b
}

func (b *CompteBuilder) SetCautionRestante(caution float64) *CompteBuilder {
	b.compte.CautionRestante = caution
	return b
}

func (b *CompteBuilder) SetIsBibliothecaire(flag bool) *CompteBuilder {
	b.compte.IsBibliothecaire = flag
	return b
}

func (b *CompteBuilder) Build() Compte {
	return b.compte
}

func (b *CompteBuilder) Reset() {
	b.compte = Compte{}
}
