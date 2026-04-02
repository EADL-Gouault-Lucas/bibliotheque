package models

type LivreBuilder struct {
	livre Livre
}

func NewLivreBuilder() *LivreBuilder {
	return &LivreBuilder{}
}

func (b *LivreBuilder) SetID(id uint) *LivreBuilder {
	b.livre.ID = id
	return b
}

func (b *LivreBuilder) SetTitre(titre string) *LivreBuilder {
	b.livre.Titre = titre
	return b
}

func (b *LivreBuilder) SetCodeBarre(codeBarre string) *LivreBuilder {
	b.livre.CodeBarre = codeBarre
	return b
}

func (b *LivreBuilder) SetCodeISBN(codeIsbn string) *LivreBuilder {
	b.livre.CodeISBN = codeIsbn
	return b
}

func (b *LivreBuilder) SetAuteurs(auteurs []string) *LivreBuilder {
	b.livre.Auteurs = auteurs
	return b
}

func (b *LivreBuilder) Build() Livre {
	return b.livre
}

func (b *LivreBuilder) Reset() {
	b.livre = Livre{}
}
