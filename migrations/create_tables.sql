DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role VARCHAR(100)
);

CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(100) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100) NOT NULL,
    surname    VARCHAR(100),
    role_id    INT REFERENCES roles(id)
);

CREATE TABLE sessions (
    id            SERIAL PRIMARY KEY,
    user_id       INT REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(300),
    expires_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

INSERT INTO roles (role) 
VALUES ('student'), ('teacher'), ('admin');

INSERT INTO users (email, password, name, patronymic, surname, role_id) 
VALUES (
    'spamer@gmail.com', '88888888', 'ivan', 'ivanovich', 'ivanov', 3
);