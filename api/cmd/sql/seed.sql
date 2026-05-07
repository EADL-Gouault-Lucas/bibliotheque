-- Jeu de données
-- Mot de passe de tous les comptes : password123

-- ── Comptes ───────────────────────────────────────────────────────────────────
INSERT INTO comptes (created_at, updated_at, email, prenom, nom, mot_de_passe, caution_restante, is_bibliothecaire)
VALUES
  (NOW(), NOW(), 'biblio@bibliotheque.fr', 'Admin',  'Biblio', 'a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4:8c0cefbe15bab070cf1f90706cd79b6f03074ca93f78b100fcb7ea06e93f1f64', 0,  TRUE),
  (NOW(), NOW(), 'alice@exemple.fr',       'Alice',  'Martin', 'b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5:49065e7565a614b2bb4812e037e3618150bb74d86a39c0f1116e6c5c856ff276', 50, FALSE),
  (NOW(), NOW(), 'bob@exemple.fr',         'Bob',    'Dupont', 'c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6:0737f2539ce3462c610aea9e9c7d5d7c8f223c783688e17064200798f0065720', 20, FALSE)
ON CONFLICT (email) DO NOTHING;

-- ── Livres ────────────────────────────────────────────────────────────────────
INSERT INTO livres (created_at, updated_at, titre, code_isbn, auteurs)
VALUES
  (NOW(), NOW(), 'The Go Programming Language',    '978-0134190440', '["Alan Donovan","Brian Kernighan"]'),
  (NOW(), NOW(), 'Domain-Driven Design',           '978-0321125217', '["Eric Evans"]'),
  (NOW(), NOW(), 'The C++ Programming Language',   '978-0321563842', '["Bjarne Stroustrup"]')
ON CONFLICT (code_isbn) DO NOTHING;

-- ── Exemplaires ───────────────────────────────────────────────────────────────
INSERT INTO exemplaires (created_at, updated_at, code_barre, caution, travee, etagere, niveau, est_emprunte, livre_id)
SELECT NOW(), NOW(), 'GO-001',  15, 'A', '1', '2', FALSE, id FROM livres WHERE code_isbn = '978-0134190440'
ON CONFLICT (code_barre) DO NOTHING;

INSERT INTO exemplaires (created_at, updated_at, code_barre, caution, travee, etagere, niveau, est_emprunte, livre_id)
SELECT NOW(), NOW(), 'GO-002',  15, 'A', '1', '3', FALSE, id FROM livres WHERE code_isbn = '978-0134190440'
ON CONFLICT (code_barre) DO NOTHING;

INSERT INTO exemplaires (created_at, updated_at, code_barre, caution, travee, etagere, niveau, est_emprunte, livre_id)
SELECT NOW(), NOW(), 'DDD-001', 20, 'B', '2', '1', FALSE, id FROM livres WHERE code_isbn = '978-0321125217'
ON CONFLICT (code_barre) DO NOTHING;

INSERT INTO exemplaires (created_at, updated_at, code_barre, caution, travee, etagere, niveau, est_emprunte, livre_id)
SELECT NOW(), NOW(), 'DDD-002', 20, 'B', '2', '2', FALSE, id FROM livres WHERE code_isbn = '978-0321125217'
ON CONFLICT (code_barre) DO NOTHING;

INSERT INTO exemplaires (created_at, updated_at, code_barre, caution, travee, etagere, niveau, est_emprunte, livre_id)
SELECT NOW(), NOW(), 'CPP-001', 10, 'C', '3', '1', FALSE, id FROM livres WHERE code_isbn = '978-0321563842'
ON CONFLICT (code_barre) DO NOTHING;

-- ── Emprunts ──────────────────────────────────────────────────────────────────
-- Alice a GO-001 (actif, dans les délais)
INSERT INTO emprunts (created_at, updated_at, date_emprunt, date_limite, rendu, compte_id, exemplaire_id)
SELECT NOW(), NOW(),
  NOW() - INTERVAL '5 days',
  NOW() + INTERVAL '9 days',
  FALSE,
  (SELECT id FROM comptes WHERE email = 'alice@exemple.fr'),
  (SELECT id FROM exemplaires WHERE code_barre = 'GO-001')
WHERE NOT EXISTS (
  SELECT 1 FROM emprunts WHERE compte_id = (SELECT id FROM comptes WHERE email = 'alice@exemple.fr')
    AND exemplaire_id = (SELECT id FROM exemplaires WHERE code_barre = 'GO-001') AND rendu = FALSE
);
UPDATE exemplaires SET est_emprunte = TRUE WHERE code_barre = 'GO-001';

-- Alice a DDD-001 (en retard)
INSERT INTO emprunts (created_at, updated_at, date_emprunt, date_limite, rendu, compte_id, exemplaire_id)
SELECT NOW(), NOW(),
  NOW() - INTERVAL '2 months',
  NOW() - INTERVAL '1 month',
  FALSE,
  (SELECT id FROM comptes WHERE email = 'alice@exemple.fr'),
  (SELECT id FROM exemplaires WHERE code_barre = 'DDD-001')
WHERE NOT EXISTS (
  SELECT 1 FROM emprunts WHERE compte_id = (SELECT id FROM comptes WHERE email = 'alice@exemple.fr')
    AND exemplaire_id = (SELECT id FROM exemplaires WHERE code_barre = 'DDD-001') AND rendu = FALSE
);
UPDATE exemplaires SET est_emprunte = TRUE WHERE code_barre = 'DDD-001';

-- Bob a rendu GO-002
INSERT INTO emprunts (created_at, updated_at, date_emprunt, date_limite, date_retour, rendu, compte_id, exemplaire_id)
SELECT NOW(), NOW(),
  NOW() - INTERVAL '1 month',
  NOW() - INTERVAL '4 days',
  NOW() - INTERVAL '3 days',
  TRUE,
  (SELECT id FROM comptes WHERE email = 'bob@exemple.fr'),
  (SELECT id FROM exemplaires WHERE code_barre = 'GO-002')
WHERE NOT EXISTS (
  SELECT 1 FROM emprunts WHERE compte_id = (SELECT id FROM comptes WHERE email = 'bob@exemple.fr')
    AND exemplaire_id = (SELECT id FROM exemplaires WHERE code_barre = 'GO-002') AND rendu = TRUE
);
