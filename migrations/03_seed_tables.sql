INSERT INTO roles (name) 
VALUES ('student'), ('teacher'), ('admin');

INSERT INTO educational_envs (name) 
VALUES ('admin'), ('voenmeh');

INSERT INTO departments (name, env_id) 
VALUES ('admin', 1), ('О7', 2), ('О6', 2), ('О4', 2), ('И9', 2), ('Р1', 2);

INSERT INTO users (email, password, name, patronymic, surname, role_id, dep_id) 
VALUES 
('spamer@gmail.com', '88888888', 'ivan', 'ivanovich', 'ivanov', 3, 1),
('teacher@mail.ru', '88888888', 'koly', 'pidor', 'fokin', 2, 2);

INSERT INTO courses (name, term, teacher_id, markdown, dep_id)
VALUES 
('Компьютерный практикум', 1, 1, 'Some markdown1', 2),
('Информационные системы и технологии', 2, 1, 'Some markdown2', 2),
('Информационные системы и технологии', 3, 1, 'Some markdown3', 2),
('Философия', 4, 1, 'Some markdown4', 6),
('Психология', 5, 1, 'Some markdown5', 6),
('Большевистская железная дорога', 6, 1, 'Some markdown6', 6);

INSERT INTO user_courses (user_id, course_id) 
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6);