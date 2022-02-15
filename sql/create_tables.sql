CREATE TABLE IF NOT EXISTS items(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--seed
INSERT INTO items (name, description) VALUES ('first product', 'this is the first product');
INSERT INTO items (name, description) VALUES ('second product', 'this is the second product');
INSERT INTO items (name, description) VALUES ('third product', 'this is the third product');
INSERT INTO items (name, description) VALUES ('fourth product', 'this is the fourth product');