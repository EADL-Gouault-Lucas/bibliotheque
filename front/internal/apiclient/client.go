package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bibliotheque-front/internal/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ── Auth ──────────────────────────────────────────────────────────────────────

type RegisterInput struct {
	Email           string  `json:"email"`
	Prenom          string  `json:"prenom"`
	Nom             string  `json:"nom"`
	MotDePasse      string  `json:"mot_de_passe"`
	CautionInitiale float64 `json:"caution_initiale"`
	Adresse         string  `json:"adresse"`
}

type LoginInput struct {
	Email      string `json:"email"`
	MotDePasse string `json:"mot_de_passe"`
}

type LoginResponse struct {
	Token  string        `json:"token"`
	Compte models.Compte `json:"compte"`
}

func (c *Client) Register(input RegisterInput) (*models.Compte, error) {
	var compte models.Compte
	return &compte, c.post("/api/v1/auth/register", "", input, &compte)
}

func (c *Client) Login(input LoginInput) (*LoginResponse, error) {
	var resp LoginResponse
	return &resp, c.post("/api/v1/auth/login", "", input, &resp)
}

// ── Livres ────────────────────────────────────────────────────────────────────

type CreateLivreInput struct {
	Titre     string   `json:"titre"`
	CodeBarre string   `json:"code_barre"`
	CodeISBN  string   `json:"code_isbn"`
	Auteurs   []string `json:"auteurs"`
}

type AddExemplaireInput struct {
	CodeBarre string  `json:"code_barre"`
	Caution   float64 `json:"caution"`
	Travee    string  `json:"travee"`
	Etagere   string  `json:"etagere"`
	Niveau    string  `json:"niveau"`
}

func (c *Client) ListLivres() ([]models.Livre, error) {
	var livres []models.Livre
	return livres, c.get("/api/v1/livres", "", &livres)
}

func (c *Client) CreateLivre(token string, input CreateLivreInput) (*models.Livre, error) {
	var livre models.Livre
	return &livre, c.post("/api/v1/livres", token, input, &livre)
}

func (c *Client) AddExemplaire(token string, livreID uint, input AddExemplaireInput) (*models.Exemplaire, error) {
	var ex models.Exemplaire
	return &ex, c.post(fmt.Sprintf("/api/v1/livres/%d/exemplaires", livreID), token, input, &ex)
}

func (c *Client) GetLivre(livreID uint) (*models.Livre, error) {
	var livre models.Livre
	return &livre, c.get(fmt.Sprintf("/api/v1/livres/%d", livreID), "", &livre)
}

// ── Emprunts ──────────────────────────────────────────────────────────────────

type CreateEmpruntInput struct {
	ExemplaireID uint `json:"exemplaire_id"`
}

type RappelsResponse struct {
	RappelsEnvoyes int             `json:"rappels_envoyes"`
	Comptes        []models.Compte `json:"comptes"`
}

func (c *Client) MesEmprunts(token string) ([]models.Emprunt, error) {
	var emprunts []models.Emprunt
	return emprunts, c.get("/api/v1/emprunts", token, &emprunts)
}

func (c *Client) CreateEmprunt(token string, input CreateEmpruntInput) (*models.Emprunt, error) {
	var emprunt models.Emprunt
	return &emprunt, c.post("/api/v1/emprunts", token, input, &emprunt)
}

func (c *Client) RetourExemplaire(token string, empruntID uint) (*models.Emprunt, error) {
	var emprunt models.Emprunt
	return &emprunt, c.put(fmt.Sprintf("/api/v1/emprunts/%d/retour", empruntID), token, nil, &emprunt)
}

func (c *Client) ListRetards(token string) ([]models.Emprunt, error) {
	var emprunts []models.Emprunt
	return emprunts, c.get("/api/v1/emprunts/retards", token, &emprunts)
}

func (c *Client) EnvoyerRappels(token string) (*RappelsResponse, error) {
	var resp RappelsResponse
	return &resp, c.post("/api/v1/emprunts/rappels", token, nil, &resp)
}

// ── HTTP helpers ──────────────────────────────────────────────────────────────

func (c *Client) get(path, token string, out interface{}) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.do(req, out)
}

func (c *Client) post(path, token string, body, out interface{}) error {
	return c.request(http.MethodPost, path, token, body, out)
}

func (c *Client) put(path, token string, body, out interface{}) error {
	return c.request(http.MethodPut, path, token, body, out)
}

func (c *Client) request(method, path, token string, body, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.do(req, out)
}

type apiError struct {
	Err string `json:"error"`
}

func (c *Client) do(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req) // #nosec G704 -- baseURL is set from trusted configuration
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var apiErr apiError
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Err != "" {
			return fmt.Errorf("%s", apiErr.Err)
		}
		return fmt.Errorf("erreur %d", resp.StatusCode)
	}

	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}
