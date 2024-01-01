ALTER TABLE rdtr_group
ADD COLUMN asset_key VARCHAR(255) NULL;
ALTER TABLE rdtr_group
DROP COLUMN cat_key;