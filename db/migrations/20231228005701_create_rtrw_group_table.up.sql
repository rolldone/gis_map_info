CREATE TABLE rtrw_group (
    id BIGINT PRIMARY KEY,
    rtrw_id BIGINT,
    properties JSONB,
    status VARCHAR(255),
    name VARCHAR(255),
    cat_key VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);