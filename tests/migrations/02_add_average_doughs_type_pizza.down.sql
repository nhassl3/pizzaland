DROP TABLE pizza_dough_types;
DROP TABLE pizza_images;

ALTER TABLE pizza ADD COLUMN type_dough INTEGER NOT NULL;

ALTER TABLE pizza DROP COLUMN created_at;