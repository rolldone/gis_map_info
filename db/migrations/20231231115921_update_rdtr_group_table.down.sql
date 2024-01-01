ALTER TABLE rdtr_group
ADD COLUMN cat_key VARCHAR(255) NULL;
ALTER TABLE rdtr_group
DROP COLUMN asset_key;