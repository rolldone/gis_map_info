ALTER TABLE zlp_group
ADD COLUMN asset_key VARCHAR(255) NULL;
ALTER TABLE zlp_group
ADD COLUMN uuid UUID UNIQUE NULL;
ALTER TABLE zlp_group
DROP COLUMN cat_key;