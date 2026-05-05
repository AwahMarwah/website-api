CREATE TABLE brands
(
    id         TEXT PRIMARY KEY,
    name       TEXT                                NOT NULL,
    slug       TEXT                                NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE TABLE categories
(
    id         TEXT PRIMARY KEY,
    name       TEXT                                NOT NULL,
    slug       TEXT                                NOT NULL UNIQUE,
    parent_id  TEXT                                REFERENCES categories (id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE TABLE products
(
    id          TEXT PRIMARY KEY,
    brand_id    TEXT                                                                             REFERENCES brands (id) ON DELETE SET NULL,
    sku         TEXT                                                                             NOT NULL UNIQUE,
    name        TEXT                                                                             NOT NULL,
    slug        TEXT                                                                             NOT NULL UNIQUE,
    description TEXT,
    base_price  NUMERIC(15, 2)                                                                   NOT NULL,
    status      TEXT CHECK (status IN ('draft', 'active', 'inactive')) DEFAULT 'draft',
    created_at  TIMESTAMP                                              DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP
);


CREATE TABLE product_variants
(
    id           TEXT PRIMARY KEY,
    product_id   TEXT                                NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    sku          TEXT                                NOT NULL UNIQUE,
    variant_name TEXT                                NOT NULL,
    price        NUMERIC(15, 2)                      NOT NULL,
    stock        INTEGER                             NOT NULL DEFAULT 0,
    weight       NUMERIC(10, 2),
    is_active    BOOLEAN   DEFAULT TRUE              NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP
);

CREATE TABLE product_images
(
    id         TEXT PRIMARY KEY,
    product_id TEXT                                NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    image_url  TEXT                                NOT NULL,
    is_primary BOOLEAN   DEFAULT FALSE             NOT NULL,
    sort_order INTEGER   DEFAULT 0                 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE product_categories
(
    product_id  TEXT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE TABLE reviews
(
    id         TEXT PRIMARY KEY,
    product_id TEXT                                   NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    user_id    TEXT                                   NOT NULL,
    rating     INTEGER CHECK (rating BETWEEN 1 AND 5) NOT NULL,
    comment    TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP    NOT NULL
);

CREATE TABLE cart_items
(
    id                 TEXT PRIMARY KEY,
    user_id            TEXT                                NOT NULL,
    product_variant_id TEXT                                NOT NULL REFERENCES product_variants (id) ON DELETE CASCADE,
    qty                INTEGER                             NOT NULL CHECK (qty > 0),
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (user_id, product_variant_id)
);

CREATE TABLE order_items
(
    id                 TEXT PRIMARY KEY,
    order_id           TEXT           NOT NULL,
    product_variant_id TEXT           NOT NULL REFERENCES product_variants (id),
    price              NUMERIC(15, 2) NOT NULL,
    qty                INTEGER        NOT NULL,
    subtotal           NUMERIC(15, 2) NOT NULL
);

CREATE INDEX idx_products_status ON products (status);
CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);
CREATE INDEX idx_product_images_product_id ON product_images (product_id);
CREATE INDEX idx_reviews_product_id ON reviews (product_id);
