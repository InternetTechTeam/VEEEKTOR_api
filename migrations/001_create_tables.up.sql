DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id        SERIAL PRIMARY KEY,
    email     VARCHAR(100) UNIQUE NOT NULL,
    password  VARCHAR(100) NOT NULL,
    name      VARCHAR(100) NOT NULL,
    surname   VARCHAR(100),
    last_name VARCHAR(100) NOT NULL
);
