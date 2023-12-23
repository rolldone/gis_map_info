CREATE TABLE rtrw (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    reg_province_id INTEGER CHECK (reg_province_id >= 0),
    reg_regency_id INTEGER CHECK (reg_regency_id >= 0),
    reg_district_id INTEGER CHECK (reg_district_id >= 0),
    reg_village_id INTEGER CHECK (reg_village_id >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);