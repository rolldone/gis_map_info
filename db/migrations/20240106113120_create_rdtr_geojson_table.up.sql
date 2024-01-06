CREATE TABLE rdtr_geojson (
    uuid UUID UNIQUE NOT NULL,
    rdtr_file_id BIGINT,
    rdtr_group_id BIGINT,
    rdtr_id BIGINT,
    geojson GEOMETRY(MultiPolygonZ, 4326) NULL, -- Using POINT and SRID 4326 for geographic coordinates
    properties JSONB NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (uuid)
);
