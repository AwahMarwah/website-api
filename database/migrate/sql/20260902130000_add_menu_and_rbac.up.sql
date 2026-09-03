CREATE TABLE IF NOT EXISTS menus
(
    id           TEXT PRIMARY KEY,
    parent_id    TEXT REFERENCES menus (id),
    name         VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(100)        NOT NULL,
    icon         VARCHAR(100),
    path         VARCHAR(200),
    sort_order   INT                 DEFAULT 0,
    is_active    BOOLEAN             DEFAULT TRUE,
    created_at   TIMESTAMP           DEFAULT CURRENT_TIMESTAMP,
    created_by   TEXT,
    updated_at   TIMESTAMP,
    updated_by   TEXT,
    deleted_at   TIMESTAMP,
    deleted_by   TEXT
);

CREATE TABLE IF NOT EXISTS role_menus
(
    id         TEXT PRIMARY KEY,
    role_id    TEXT      NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    menu_id    TEXT      NOT NULL REFERENCES menus (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (role_id, menu_id)
);

ALTER TABLE permissions ADD COLUMN IF NOT EXISTS menu_id TEXT REFERENCES menus(id);

CREATE INDEX IF NOT EXISTS idx_permissions_menu_id ON permissions (menu_id);
CREATE INDEX IF NOT EXISTS idx_role_menus_role_id ON role_menus (role_id);
