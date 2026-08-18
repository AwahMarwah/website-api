ALTER TABLE user_addresses
ADD COLUMN province_id VARCHAR(2),
ADD COLUMN city_id VARCHAR(4),
ADD COLUMN district_id VARCHAR(7),
ADD COLUMN subdistrict_id VARCHAR(10),
ADD COLUMN destination_id BIGINT;

CREATE INDEX idx_user_address_province_id
    ON user_addresses(province_id);

CREATE INDEX idx_user_address_city_id
    ON user_addresses(city_id);

CREATE INDEX idx_user_address_destination_id
    ON user_addresses(destination_id);

CREATE INDEX idx_user_address_user_id
    ON user_addresses(user_id);