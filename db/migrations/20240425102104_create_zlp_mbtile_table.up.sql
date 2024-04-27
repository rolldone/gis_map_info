CREATE TABLE zlp_mbtile (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NULL,
    file_name VARCHAR(255) NULL,
    zlp_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checked_at TIMESTAMP NULL
);