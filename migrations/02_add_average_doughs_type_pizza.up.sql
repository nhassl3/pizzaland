CREATE TABLE IF NOT EXISTS pizza_temp (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     category_id INTEGER NOT NULL,
                                     name VARCHAR(50) NOT NULL UNIQUE,
                                     description VARCHAR(256),
                                     price REAL NOT NULL DEFAULT 109,
                                     diameter INTEGER,
                                     CHECK (diameter IN (26, 30, 40)),
                                     FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

INSERT INTO pizza_temp (id, category_id, name, description, price, diameter)
SELECT id, category_id, name, description, price, diameter FROM pizza;

DROP TABLE pizza;
ALTER TABLE pizza_temp RENAME TO pizza;

CREATE INDEX IF NOT EXISTS idx_pizza_name ON pizza(name);

CREATE TABLE IF NOT EXISTS pizza_dough_types (
                                                 pizza_id INTEGER REFERENCES pizza(id) ON DELETE CASCADE,
                                                 dough_type_id INTEGER REFERENCES doughs(id) ON DELETE CASCADE,
                                                 PRIMARY KEY (pizza_id, dough_type_id)
);
