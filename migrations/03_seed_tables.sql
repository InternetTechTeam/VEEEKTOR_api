INSERT INTO roles (name) 
VALUES ('student'), ('teacher'), ('admin');

INSERT INTO educational_envs (name) 
VALUES ('admin'), ('voenmeh');

INSERT INTO departments (name, env_id) 
VALUES ('admin', 1), ('О7', 2), ('О6', 2), ('О4', 2), ('И9', 2);

INSERT INTO users (email, password, name, patronymic, surname, role_id, dep_id) 
VALUES 
('spamer@gmail.com', '88888888', 'ivan', 'ivanovich', 'ivanov', 3, 1);