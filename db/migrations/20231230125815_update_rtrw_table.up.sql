ALTER TABLE rtrw
ADD COLUMN status VARCHAR(20) NULL;

ALTER TABLE rtrw
ALTER COLUMN reg_province_id TYPE BIGINT USING reg_province_id::BIGINT;
ALTER TABLE rtrw
ALTER COLUMN reg_regency_id TYPE BIGINT USING reg_regency_id::BIGINT;
ALTER TABLE rtrw
ALTER COLUMN reg_district_id TYPE BIGINT USING reg_district_id::BIGINT;
ALTER TABLE rtrw
ALTER COLUMN reg_village_id TYPE BIGINT USING reg_village_id::BIGINT;