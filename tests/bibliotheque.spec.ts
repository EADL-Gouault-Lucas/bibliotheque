import { test, expect } from '@playwright/test';

// Comptes disponibles via make sql (mot de passe : password123)
const ALICE = { email: 'alice@exemple.fr', password: 'password123' };
const BIBLIO = { email: 'biblio@bibliotheque.fr', password: 'password123' };

// ── Helpers ────────────────────────────────────────────────────────────────────

async function login(page: any, email: string, password: string) {
    await page.goto('/login');
    await page.fill('#email', email);
    await page.fill('#mot_de_passe', password);
    await page.click('button[type="submit"]');
}

// ── Catalogue (public) ─────────────────────────────────────────────────────────

test.describe('Catalogue', () => {
    test('affiche la page catalogue sans connexion', async ({ page }) => {
        await page.goto('/');
        await expect(page).toHaveTitle(/Catalogue/);
        await expect(page.getByRole('heading', { name: 'Catalogue des livres' })).toBeVisible();
    });

    test('liste les livres du catalogue', async ({ page }) => {
        await page.goto('/');
        await expect(page.getByText('The Go Programming Language')).toBeVisible();
        await expect(page.getByText('Domain-Driven Design')).toBeVisible();
    });

    test('n\'affiche pas le bouton "Ajouter un livre" sans connexion', async ({ page }) => {
        await page.goto('/');
        await expect(page.getByRole('link', { name: '+ Ajouter un livre' })).not.toBeVisible();
    });
});

// ── Authentification ───────────────────────────────────────────────────────────

test.describe('Connexion', () => {
    test('affiche la page de connexion', async ({ page }) => {
        await page.goto('/login');
        await expect(page).toHaveTitle(/Connexion/);
        await expect(page.getByRole('button', { name: 'Se connecter' })).toBeVisible();
    });

    test('connexion réussie redirige vers le catalogue', async ({ page }) => {
        await login(page, ALICE.email, ALICE.password);
        await expect(page).toHaveURL('/');
        await expect(page.getByRole('heading', { name: 'Catalogue des livres' })).toBeVisible();
    });

    test('connexion échoue avec un mauvais mot de passe', async ({ page }) => {
        await login(page, ALICE.email, 'mauvais_mot_de_passe');
        await expect(page).toHaveURL(/\/login/);
        await expect(page.locator('.alert-danger')).toBeVisible();
    });

    test('lien vers inscription présent sur la page de connexion', async ({ page }) => {
        await page.goto('/login');
        await expect(page.getByRole('link', { name: "S'inscrire" })).toBeVisible();
    });
});

// ── Inscription ────────────────────────────────────────────────────────────────

test.describe('Inscription', () => {
    test('affiche la page d\'inscription', async ({ page }) => {
        await page.goto('/register');
        await expect(page).toHaveTitle(/Inscription/);
        await expect(page.locator('form[action="/register"]')).toBeVisible();
    });

    test('inscription réussie redirige vers la connexion', async ({ page }) => {
        const unique = Date.now();
        await page.goto('/register');
        await page.fill('#prenom', 'Test');
        await page.fill('#nom', 'Utilisateur');
        await page.fill('#email', `test${unique}@exemple.fr`);
        await page.fill('#mot_de_passe', 'password123');
        await page.click('button[type="submit"]');
        await expect(page).toHaveURL(/\/login/);
        await expect(page.locator('.alert-success')).toBeVisible();
    });

    test('inscription échoue si email déjà utilisé', async ({ page }) => {
        await page.goto('/register');
        await page.fill('#prenom', 'Alice');
        await page.fill('#nom', 'Martin');
        await page.fill('#email', ALICE.email);
        await page.fill('#mot_de_passe', 'password123');
        await page.click('button[type="submit"]');
        await expect(page).toHaveURL(/\/register/);
        await expect(page.locator('.alert-danger')).toBeVisible();
    });
});

// ── Emprunts (utilisateur connecté) ───────────────────────────────────────────

test.describe('Emprunts', () => {
    test('redirige vers /login si non connecté', async ({ page }) => {
        await page.goto('/emprunts');
        await expect(page).toHaveURL(/\/login/);
    });

    test('affiche la page emprunts après connexion', async ({ page }) => {
        await login(page, ALICE.email, ALICE.password);
        await page.goto('/emprunts');
        await expect(page).toHaveTitle(/emprunts/i);
        await expect(page.getByRole('heading', { name: 'Mes emprunts' })).toBeVisible();
    });

    test('affiche les emprunts d\'Alice', async ({ page }) => {
        await login(page, ALICE.email, ALICE.password);
        await page.goto('/emprunts');
        // Alice a GO-001 dans le seed
        await expect(page.getByText('GO-001')).toBeVisible();
    });
});

// ── Administration (bibliothécaire) ───────────────────────────────────────────

test.describe('Administration', () => {
    test('affiche le bouton "Ajouter un livre" pour le bibliothécaire', async ({ page }) => {
        await login(page, BIBLIO.email, BIBLIO.password);
        await page.goto('/');
        await expect(page.getByRole('link', { name: '+ Ajouter un livre' })).toBeVisible();
    });

    test('accès à la page de création de livre', async ({ page }) => {
        await login(page, BIBLIO.email, BIBLIO.password);
        await page.goto('/admin/livres/nouveau');
        await expect(page).toHaveTitle(/livre/i);
        await expect(page.locator('form')).toBeVisible();
    });

    test('/admin/livres/nouveau redirige vers /login pour un utilisateur non connecté', async ({ page }) => {
        await page.goto('/admin/livres/nouveau');
        await expect(page).toHaveURL(/\/login/);
    });

    test('/admin/livres/nouveau redirige pour un utilisateur non bibliothécaire', async ({ page }) => {
        await login(page, ALICE.email, ALICE.password);
        await page.goto('/admin/livres/nouveau');
        // Alice n'est pas bibliothécaire : doit être redirigée
        await expect(page).not.toHaveURL('/admin/livres/nouveau');
    });
});

// ── Déconnexion ────────────────────────────────────────────────────────────────

test.describe('Déconnexion', () => {
    test('déconnexion redirige vers le catalogue', async ({ page }) => {
        await login(page, ALICE.email, ALICE.password);
        await page.goto('/logout');
        await expect(page).toHaveURL('/');
    });

    test('après déconnexion, /emprunts redirige vers /login', async ({ page }) => {
        await login(page, ALICE.email, ALICE.password);
        await page.goto('/logout');
        await page.goto('/emprunts');
        await expect(page).toHaveURL(/\/login/);
    });
});
