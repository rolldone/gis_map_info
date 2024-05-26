CREATE TABLE user_data (
    id bigserial PRIMARY KEY,
    uuid UUID NULL UNIQUE,
    name VARCHAR(255) NULL,
    username VARCHAR(255) NULL UNIQUE,
    email VARCHAR(255) NULL UNIQUE,
    passkey TEXT NULL,
    salt TEXT NULL,
    status VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);