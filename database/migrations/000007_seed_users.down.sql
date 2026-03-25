-- Rollback seed: remove default users
DELETE FROM users WHERE id IN ('fadel', 'nanang', 'try', 'daffa', 'dani');
