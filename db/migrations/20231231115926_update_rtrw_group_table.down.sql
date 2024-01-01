ALTER TABLE rtrw_group
ADD COLUMN cat_key VARCHAR(255) NULL;
ALTER TABLE rtrw_group
DROP COLUMN asset_key;
