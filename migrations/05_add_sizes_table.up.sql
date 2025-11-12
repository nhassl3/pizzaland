CREATE TABLE IF NOT EXISTS pizza_new (
                                              id INTEGER PRIMARY KEY AUTOINCREMENT,
                                              category_id INTEGER NOT NULL,
                                              name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(256),
    price REAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    rating INTEGER,
    image_path VARCHAR(500) NOT NULL,
    CHECK (price > 109),
    CHECK (rating >= 0 AND rating <= 5),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
    );

INSERT INTO pizza_new (id, category_id, name, description, price, created_at)
SELECT id, category_id, name, description, price, created_at
FROM pizza;

DROP TABLE pizza;

ALTER TABLE pizza_new RENAME TO pizza;

CREATE TABLE IF NOT EXISTS pizza_sizes (
    pizza_id INTEGER NOT NULL REFERENCES pizza(id) ON DELETE CASCADE,
    sizes INTEGER NOT NULL,
    CHECK (sizes IN (26, 30, 40)),
    PRIMARY KEY (pizza_id, sizes)
);
