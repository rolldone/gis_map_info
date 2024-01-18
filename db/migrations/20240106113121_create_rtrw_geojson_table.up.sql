CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE rtrw_geojson (
    uuid UUID UNIQUE NOT NULL,
    order_number BIGINT NULL,
    rtrw_file_id BIGINT,
    rtrw_group_id BIGINT,
    rtrw_id BIGINT,
    geojson GEOMETRY(Polygon, 4326) NULL, 
    -- Using POINT and SRID 4326 for geographic coordinates
    -- Using POINT and SRID 3857 for web mercator
    properties JSONB NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (uuid)
);

-- Create index for order_number
CREATE INDEX idx_order_number_rtrw_geojson ON rtrw_geojson (order_number);

-- Create index for spatial geojson
CREATE INDEX idx_rtrw_geojson_geojson_gist
ON public.rtrw_geojson 
USING gist(geojson);