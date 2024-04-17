INSERT INTO roles (name) 
VALUES 
('student'), ('teacher'), ('admin');

INSERT INTO educational_envs (name) 
VALUES 
('admin'), ('voenmeh');

INSERT INTO departments (name, env_id) 
VALUES 
('admin', 1), ('О7', 2), ('О6', 2), 
('О4', 2), ('И9', 2), ('Р1', 2);

INSERT INTO users (email, password, name, patronymic, surname, role_id, dep_id) 
VALUES 
('spamer@mail.ru', '88888888', 'ivan', 'ivanovich', 'ivanov', 3, 1),
('teacher@mail.ru', '88888888', 'koly', 'pidor', 'fokin', 2, 2),
('studentO7@mail.ru', '88888888', 'anna', 'lokiv', 'bobsova', 1, 2),
('studentO6@mail.ru', '88888888', 'alex', 'mashinov', 'bobrov', 1, 3),
('studentO4@mail.ru', '88888888', 'sasha', 'teapet', 'ruric', 1, 4),
('studentI9@mail.ru', '88888888', 'vitya', 'nextov', 'kuropyat', 1, 5),
('studentP1@mail.ru', '88888888', 'maria', 'mariovna', 'petrova', 1, 6),
('studentALL@mail.ru', '88888888', 'genius', 'vse', 'kursi', 1, 2),
('teacherO6@mail.ru', '88888888', 'teacher', 'teacher', 'teacher', 2, 3),
('teacherO4@mail.ru', '88888888', 'teacher', 'teacher', 'teacher', 2, 4),
('teacherI9@mail.ru', '88888888', 'teacher', 'teacher', 'teacher', 2, 5),
('teacherP1@mail.ru', '88888888', 'teacher', 'teacher', 'teacher', 2, 6);

INSERT INTO courses (name, term, teacher_id, markdown, dep_id) 
VALUES 
('Компьютерный практикум', 1, 1, 'Some markdown1', 2),
('Информационные системы и технологии', 2, 1, 'Some markdown2', 2),
('Информационные системы и технологии', 3, 1, 'Some markdown3', 2),
('Философия', 4, 1, 'Some markdown4', 6),
('Психология', 5, 1, 'Some markdown5', 6),
('Большевистская железная дорога', 6, 1, 'Some markdown6', 6),
('Крутой предмет О6', 2, 2, '## О6, **FATO6** __underline__, > O6?', 3),
('Крутой предмет О4', 2, 3, '## О4, **FATO4** __underline__, > O4?', 4),
('Крутой предмет И9', 2, 4, '## И9, **FATИ9** __underline__, > И9?', 5),
('Крутой предмет П1', 2, 5, '## P1, **FATP1** __underline__, > P1?', 6);


INSERT INTO user_courses (user_id, course_id) 
VALUES 
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), 
(1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
(3, 1), (3, 2), (3, 3), (4, 7), (5, 8), 
(6, 9), (7, 10), (8, 1), (8, 2), (8, 3), 
(8, 4), (8, 5), (8, 6), (8, 7), (8, 8), 
(8, 9), (8, 10);

INSERT INTO locations (location) 
VALUES 
('At home'), ('In class');

INSERT INTO nested_infos (course_id, name, markdown)
VALUES 
(6, 'БЖД', 'БЖД = круто.');

INSERT INTO nested_tests 
(course_id, opens, closes, 
tasks_count, topic, location_id, 
attempts, password, time_limit)
VALUES 
(6, '1999-01-08 04:05:06', '2026-09-08 00:00:00',
15, 'Железная дорога', 2, 1, '', '00:15:00'),
(4, '1999-01-08 04:05:06', '2026-09-08 00:00:00',
15, 'Карл Маркс', 2, 1, '', '00:20:00');

INSERT INTO nested_labs 
(course_id, opens, closes, topic, requirements, 
example, location_id, attempts)
VALUES
(6, '1999-01-08 04:05:06', '2026-09-08 00:00:00',
'Лабораторная работа 1', '', '', 1, 0);
