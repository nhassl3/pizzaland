DROP TABLE IF EXISTS pizza_dough_types;
DROP TABLE IF EXISTS pizza_images;

ALTER TABLE pizza ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE TABLE IF NOT EXISTS pizza_dough_types (
    pizza_id INTEGER REFERENCES pizza(id) ON DELETE CASCADE,
    dough_type_id INTEGER REFERENCES doughs(id) ON DELETE CASCADE,
    PRIMARY KEY(pizza_id, dough_type_id)
);

CREATE TABLE IF NOT EXISTS pizza_images (
                                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                                            pizza_id INTEGER REFERENCES pizza(id) ON DELETE CASCADE,
    image_path VARCHAR(500) NOT NULL, -- or url for cloud storage
    image_type VARCHAR(20) NOT NULL, -- 'original', 'thumbnail', 'medium'
    is_main BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );