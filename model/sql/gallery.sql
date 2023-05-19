CREATE TABLE galleries (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    user INT REFERENCE users (id) NOT NULL,
);