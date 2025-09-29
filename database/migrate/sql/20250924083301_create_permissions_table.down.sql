-- Hapus tabel role_permissions terlebih dahulu (karena ada foreign key ke permissions)
DROP TABLE IF EXISTS role_permissions;

-- Lalu hapus tabel permissions
DROP TABLE IF EXISTS permissions;