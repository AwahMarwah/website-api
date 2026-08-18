CREATE TABLE provinces
(
    id         VARCHAR(2) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE cities
(
    id          VARCHAR(4) PRIMARY KEY,
    province_id VARCHAR(2)   NOT NULL,
    name        VARCHAR(255) NOT NULL,

    created_at  TIMESTAMP,
    updated_at  TIMESTAMP,

    CONSTRAINT fk_city_province
        FOREIGN KEY (province_id)
            REFERENCES provinces (id)
);

CREATE TABLE districts
(
    id         VARCHAR(7) PRIMARY KEY,
    city_id    VARCHAR(4)   NOT NULL,
    name       VARCHAR(255) NOT NULL,

    created_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_district_city
        FOREIGN KEY (city_id)
            REFERENCES cities (id)
);

CREATE TABLE subdistricts
(
    id          VARCHAR(10) PRIMARY KEY,
    district_id VARCHAR(7)   NOT NULL,
    name        VARCHAR(255) NOT NULL,

    created_at  TIMESTAMP,
    updated_at  TIMESTAMP,

    CONSTRAINT fk_subdistrict_district
        FOREIGN KEY (district_id)
            REFERENCES districts (id)
);

CREATE TABLE city_shipping_mapping
(
    city_id        VARCHAR(4)  NOT NULL,
    provider       VARCHAR(50) NOT NULL,
    destination_id BIGINT      NOT NULL,

    created_at     TIMESTAMP,

    PRIMARY KEY (city_id, provider)
);
