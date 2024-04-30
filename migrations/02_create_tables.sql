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
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(200) NOT NULL,
    term           INT NOT NULL,
    teacher_id     INT REFERENCES users(id) ON DELETE SET NULL,
    markdown       TEXT,
    dep_id         INT REFERENCES departments(id) ON DELETE SET NULL,
    modified_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE user_courses (
    id        SERIAL PRIMARY KEY,
    user_id   INT REFERENCES users(id) ON DELETE CASCADE,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE
);

CREATE TABLE nested_infos (
    id        SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    name      VARCHAR(512),
    markdown  TEXT
);

CREATE TABLE locations (
    id       SERIAL PRIMARY KEY,
    location VARCHAR(512)
);

CREATE TABLE nested_tests (
    id          SERIAL PRIMARY KEY, 
    course_id   INT REFERENCES courses(id) ON DELETE CASCADE,
    opens       TIMESTAMP WITH TIME ZONE NOT NULL,
    closes      TIMESTAMP WITH TIME ZONE NOT NULL,
    tasks_count INT NOT NULL,
    topic       VARCHAR(512) NOT NULL,
    location_id INT REFERENCES locations(id) ON DELETE SET NULL,
    attempts    INT NOT NULL, 
    password    VARCHAR(256),
    time_limit  TIME NOT NULL
);

CREATE TABLE nested_labs (
    id           SERIAL PRIMARY KEY,
    course_id    INT REFERENCES courses(id) ON DELETE CASCADE,
    opens        TIMESTAMP WITH TIME ZONE NOT NULL,
    closes       TIMESTAMP WITH TIME ZONE NOT NULL,
    topic        VARCHAR(512),
    requirements VARCHAR(512),
    example      VARCHAR(512),
    location_id  INT REFERENCES locations(id) ON DELETE SET NULL,
    attempts     INT NOT NULL
);
