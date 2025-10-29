CREATE TABLE IF NOT EXISTS pizza (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     category_id INTEGER NOT NULL,
                                     name VARCHAR(50) NOT NULL,
                                     description VARCHAR(256),
                                     type_dough INTEGER NOT NULL,
                                     price REAL NOT NULL DEFAULT 109,
                                     diameter INTEGER,
                                     CHECK (diameter IN (26, 30, 40)),
                                     FOREIGN KEY (type_dough) REFERENCES doughs(id),
                                     FOREIGN KEY (category_id) REFERENCES categories(id)
);
CREATE INDEX IF NOT EXISTS idx_pizza_name ON pizza(name);

CREATE TABLE IF NOT EXISTS categories (
                                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                                          name VARCHAR(26) NOT NULL,
                                          description VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS doughs (
                                      id INTEGER PRIMARY KEY AUTOINCREMENT,
                                      name VARCHAR(100) NOT NULL
);