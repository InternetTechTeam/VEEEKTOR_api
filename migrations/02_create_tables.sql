CREATE TABLE roles (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE educational_envs (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL
);

CREATE TABLE departments (
    id     SERIAL PRIMARY KEY,
    name   VARCHAR(200) NOT NULL,
    env_id INT REFERENCES educational_envs(id) ON DELETE SET NULL
);

CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(100) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100) NOT NULL,
    surname    VARCHAR(100),
    role_id    INT REFERENCES roles(id),
    dep_id     INT REFERENCES departments(id) ON DELETE SET NULL
);

CREATE TABLE sessions (
    id            SERIAL PRIMARY KEY,
    user_id       INT REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(300),
    expires_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE courses (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(200) NOT NULL,
    term       INT,
    teacher_id INT REFERENCES users(id) ON DELETE SET NULL,
    markdown   TEXT,
    dep_id     INT REFERENCES departments(id) ON DELETE SET NULL
);

CREATE TABLE user_courses (
    id        SERIAL PRIMARY KEY,
    user_id   INT REFERENCES users(id) ON DELETE CASCADE,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE
);