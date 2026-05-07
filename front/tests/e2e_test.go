//go:build e2e

package tests_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"bibliotheque-front/internal/apiclient"
	"bibliotheque-front/internal/handlers"
	"bibliotheque-front/internal/session"

	"github.com/gin-gonic/gin"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func apiURL() string {
	if v := os.Getenv("API_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

// buildRouter construit le routeur gin identique à cmd/server/main.go.
func buildRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	api := apiclient.New(apiURL())

	authHandler := handlers.NewAuthHandler(api)
	catalogueHandler := handlers.NewCatalogueHandler(api)
	empruntHandler := handlers.NewEmpruntHandler(api)
	adminHandler := handlers.NewAdminHandler(api)

	r := gin.New()

	r.GET("/", catalogueHandler.ShowCatalogue)
	r.GET("/login", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/register", authHandler.ShowRegister)
	r.POST("/register", authHandler.Register)
	r.GET("/logout", authHandler.Logout)

	auth := r.Group("/", handlers.RequireAuth())
	{
		auth.GET("/emprunts", empruntHandler.ShowMesEmprunts)
		auth.POST("/emprunts", empruntHandler.CreateEmprunt)
	}

	biblio := r.Group("/", handlers.RequireAuth(), handlers.RequireBibliothecaire())
	{
		biblio.GET("/admin/livres/nouveau", catalogueHandler.ShowNouveauLivre)
		biblio.POST("/admin/livres", catalogueHandler.CreateLivre)
		biblio.GET("/admin/livres/:id/exemplaires/nouveau", catalogueHandler.ShowNouvelExemplaire)
		biblio.POST("/admin/livres/:id/exemplaires", catalogueHandler.CreateExemplaire)
		biblio.GET("/admin/emprunts", adminHandler.ShowEmprunts)
		biblio.POST("/admin/emprunts/:id/retour", adminHandler.RetourExemplaire)
		biblio.POST("/admin/emprunts/rappels", adminHandler.EnvoyerRappels)
	}

	return r
}

// newClient crée un client HTTP qui ne suit pas les redirections.
func newClientNoRedirect() *http.Client {
	return &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// newClientWithJar crée un client HTTP avec jar de cookies (suit les redirections).
func newClientWithJar(baseURL string) *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			// Réécrire les redirections vers le serveur de test
			req.URL.Host = strings.TrimPrefix(baseURL, "http://")
			return nil
		},
	}
	return client
}

// makeSessionCookie encode manuellement un User en cookie de session.
func makeSessionCookie(user session.User) *http.Cookie {
	data, _ := json.Marshal(user)
	encoded := base64.URLEncoding.EncodeToString(data)
	return &http.Cookie{Name: "bib_session", Value: encoded}
}

// apiLoginDirect appelle directement l'API pour obtenir un token.
func apiLoginDirect(t *testing.T, email, password string) (token string, compte apiclient.LoginResponse) {
	t.Helper()
	client := apiclient.New(apiURL())
	resp, err := client.Login(apiclient.LoginInput{Email: email, MotDePasse: password})
	if err != nil {
		t.Fatalf("login API échoué pour %s : %v", email, err)
	}
	return resp.Token, *resp
}

func body(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("lecture body : %v", err)
	}
	return string(b)
}

// ── Tests : pages publiques ───────────────────────────────────────────────────

func TestCatalogue_AfficheLesCatalogue(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET / : %v", err)
	}
	b := body(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("statut attendu 200, obtenu %d", resp.StatusCode)
	}
	if !strings.Contains(b, "Catalogue") {
		t.Error("la page ne contient pas 'Catalogue'")
	}
}

func TestCatalogue_ContientLivresDuSeed(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET / : %v", err)
	}
	b := body(t, resp)

	if !strings.Contains(b, "The Go Programming Language") {
		t.Error("le catalogue ne contient pas 'The Go Programming Language'")
	}
	if !strings.Contains(b, "Domain-Driven Design") {
		t.Error("le catalogue ne contient pas 'Domain-Driven Design'")
	}
}

func TestLogin_AfficheLaPage(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/login")
	if err != nil {
		t.Fatalf("GET /login : %v", err)
	}
	b := body(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("statut attendu 200, obtenu %d", resp.StatusCode)
	}
	if !strings.Contains(b, "Connexion") {
		t.Error("la page ne contient pas 'Connexion'")
	}
}

func TestRegister_AfficheLaPage(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/register")
	if err != nil {
		t.Fatalf("GET /register : %v", err)
	}
	b := body(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("statut attendu 200, obtenu %d", resp.StatusCode)
	}
	if !strings.Contains(b, "Creer un compte") {
		t.Error("la page ne contient pas 'Creer un compte'")
	}
}

// ── Tests : authentification ──────────────────────────────────────────────────

func TestLogin_Succes_PoseCookie(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	client := newClientNoRedirect()
	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"email":        {"alice@exemple.fr"},
		"mot_de_passe": {"password123"},
	})
	if err != nil {
		t.Fatalf("POST /login : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("statut attendu 302, obtenu %d", resp.StatusCode)
	}

	var hasCookie bool
	for _, c := range resp.Cookies() {
		if c.Name == "bib_session" {
			hasCookie = true
			break
		}
	}
	if !hasCookie {
		t.Error("cookie bib_session absent après connexion réussie")
	}

	if loc := resp.Header.Get("Location"); loc != "/" {
		t.Errorf("redirection attendue vers '/', obtenu '%s'", loc)
	}
}

func TestLogin_Echec_RedirigeSurLogin(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	client := newClientNoRedirect()
	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"email":        {"alice@exemple.fr"},
		"mot_de_passe": {"mauvais_mot_de_passe"},
	})
	if err != nil {
		t.Fatalf("POST /login : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("statut attendu 302, obtenu %d", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); !strings.HasPrefix(loc, "/login") {
		t.Errorf("redirection attendue vers /login, obtenu '%s'", loc)
	}
}

func TestLogout_VideLaSession(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	client := newClientNoRedirect()
	resp, err := client.Get(srv.URL + "/logout")
	if err != nil {
		t.Fatalf("GET /logout : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("statut attendu 302, obtenu %d", resp.StatusCode)
	}
	// Le cookie doit être supprimé (MaxAge=-1)
	for _, c := range resp.Cookies() {
		if c.Name == "bib_session" && c.MaxAge > 0 {
			t.Error("le cookie bib_session devrait être supprimé après logout")
		}
	}
}

// ── Tests : protection par authentification ───────────────────────────────────

func TestEmprunts_SansAuth_RedirigeSurLogin(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	client := newClientNoRedirect()
	resp, err := client.Get(srv.URL + "/emprunts")
	if err != nil {
		t.Fatalf("GET /emprunts : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("statut attendu 302, obtenu %d", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc != "/login" {
		t.Errorf("redirection attendue vers /login, obtenu '%s'", loc)
	}
}

func TestEmprunts_AvecAuth_AfficheLaPage(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	token, resp := apiLoginDirect(t, "alice@exemple.fr", "password123")
	cookie := makeSessionCookie(session.User{
		CompteID:         resp.Compte.ID,
		Prenom:           resp.Compte.Prenom,
		Nom:              resp.Compte.Nom,
		Token:            token,
		IsBibliothecaire: resp.Compte.IsBibliothecaire,
	})

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/emprunts", nil)
	req.AddCookie(cookie)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET /emprunts : %v", err)
	}
	b := body(t, httpResp)

	if httpResp.StatusCode != http.StatusOK {
		t.Errorf("statut attendu 200, obtenu %d", httpResp.StatusCode)
	}
	if !strings.Contains(b, "Mes emprunts") {
		t.Error("la page ne contient pas 'Mes emprunts'")
	}
}

// ── Tests : protection bibliothécaire ─────────────────────────────────────────

func TestAdminLivres_SansAuth_RedirigeSurLogin(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	client := newClientNoRedirect()
	resp, err := client.Get(srv.URL + "/admin/livres/nouveau")
	if err != nil {
		t.Fatalf("GET /admin/livres/nouveau : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("statut attendu 302, obtenu %d", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc != "/login" {
		t.Errorf("redirection attendue vers /login, obtenu '%s'", loc)
	}
}

func TestAdminLivres_UtilisateurNormal_Redirige(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	token, resp := apiLoginDirect(t, "alice@exemple.fr", "password123")
	cookie := makeSessionCookie(session.User{
		CompteID:         resp.Compte.ID,
		Prenom:           resp.Compte.Prenom,
		Nom:              resp.Compte.Nom,
		Token:            token,
		IsBibliothecaire: resp.Compte.IsBibliothecaire,
	})

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/admin/livres/nouveau", nil)
	req.AddCookie(cookie)

	client := newClientNoRedirect()
	httpResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET /admin/livres/nouveau : %v", err)
	}
	defer httpResp.Body.Close()

	// Alice n'est pas bibliothécaire → redirigée vers /
	if httpResp.StatusCode != http.StatusFound {
		t.Errorf("statut attendu 302, obtenu %d", httpResp.StatusCode)
	}
	if loc := httpResp.Header.Get("Location"); loc != "/" {
		t.Errorf("redirection attendue vers '/', obtenu '%s'", loc)
	}
}

func TestAdminLivres_Bibliothecaire_AfficheLaPage(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	token, resp := apiLoginDirect(t, "biblio@bibliotheque.fr", "password123")
	cookie := makeSessionCookie(session.User{
		CompteID:         resp.Compte.ID,
		Prenom:           resp.Compte.Prenom,
		Nom:              resp.Compte.Nom,
		Token:            token,
		IsBibliothecaire: resp.Compte.IsBibliothecaire,
	})

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/admin/livres/nouveau", nil)
	req.AddCookie(cookie)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET /admin/livres/nouveau : %v", err)
	}
	b := body(t, httpResp)

	if httpResp.StatusCode != http.StatusOK {
		t.Errorf("statut attendu 200, obtenu %d", httpResp.StatusCode)
	}
	if !strings.Contains(b, "livre") {
		t.Error("la page ne contient pas 'livre'")
	}
}

// ── Tests : flux complet login → page protégée ────────────────────────────────

func TestFlux_LoginPuisEmprunts(t *testing.T) {
	srv := httptest.NewServer(buildRouter())
	defer srv.Close()

	// 1. Login via POST /login → récupère le cookie
	loginResp, err := newClientNoRedirect().PostForm(srv.URL+"/login", url.Values{
		"email":        {"alice@exemple.fr"},
		"mot_de_passe": {"password123"},
	})
	if err != nil {
		t.Fatalf("POST /login : %v", err)
	}
	defer loginResp.Body.Close()

	var sessionCookieVal string
	for _, c := range loginResp.Cookies() {
		if c.Name == "bib_session" {
			sessionCookieVal = c.Value
		}
	}
	if sessionCookieVal == "" {
		t.Fatal("cookie bib_session absent après login")
	}

	// 2. GET /emprunts avec le cookie obtenu
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/emprunts", nil)
	req.AddCookie(&http.Cookie{Name: "bib_session", Value: sessionCookieVal})

	empruntsResp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("GET /emprunts : %v", err)
	}
	b := body(t, empruntsResp)

	if empruntsResp.StatusCode != http.StatusOK {
		t.Errorf("statut attendu 200, obtenu %d", empruntsResp.StatusCode)
	}
	if !strings.Contains(b, "Mes emprunts") {
		t.Error("la page /emprunts ne contient pas 'Mes emprunts'")
	}
}
