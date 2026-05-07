-- Vide toutes les tables en respectant les contraintes de clés étrangères
TRUNCATE TABLE emprunts, exemplaires, livres, comptes RESTART IDENTITY CASCADE;
