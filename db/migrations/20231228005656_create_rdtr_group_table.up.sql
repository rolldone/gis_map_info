CREATE TABLE rdtr_group (
    id BIGSERIAL PRIMARY KEY,
    rdtr_id BIGINT,
    properties JSONB,
    status VARCHAR(255),
    name VARCHAR(255),
    cat_key VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);