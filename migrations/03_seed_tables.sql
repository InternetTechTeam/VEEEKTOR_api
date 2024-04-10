INSERT INTO roles (name) 
VALUES ('student'), ('teacher'), ('admin');

INSERT INTO educational_envs (name) 
VALUES ('admin'), ('voenmeh');

INSERT INTO departments (name, env_id) 
VALUES ('admin', 1), ('О7', 2), ('О6', 2), ('О4', 2), ('И9', 2), ('Р1', 2);

INSERT INTO users (email, password, name, patronymic, surname, role_id, dep_id) 
VALUES 
('spamer@gmail.com', '88888888', 'ivan', 'ivanovich', 'ivanov', 3, 1),
('teacher@mail.ru', '88888888', 'koly', 'pidor', 'fokin', 2, 2),
('studentO7@gmail.com', '88888888', 'anna', 'lokiv', 'bobsova', 1, 2),
('studentO6@gmail.com', '88888888', 'alex', 'mashinov', 'bobrov', 1, 3),
('studentO4@gmail.com', '88888888', 'sasha', 'teapet', 'ruric', 1, 4),
('studentI9@gmail.com', '88888888', 'vitya', 'nextov', 'kuropyat', 1, 5),
('studentP1@gmail.com', '88888888', 'maria', 'mariovna', 'petrova', 1, 6),
('studentP1@gmail.com', '88888888', 'maria', 'mariovna', 'petrova', 1, 6),
('studentALL@gmail.com', '88888888', 'genius', 'vse', 'kursi', 1, 2),
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
('Крутой предмет О6', 2, 2, '##О6, **FATO6** __underline__, >O6?', 3),
('Крутой предмет О4', 2, 3, '##О4, **FATO4** __underline__, >O4?', 4),
('Крутой предмет И9', 2, 4, '##И9, **FATИ9** __underline__, >И9?', 5),
('Крутой предмет П1', 2, 5, '##P1, **FATP1** __underline__, >P1?', 6);


INSERT INTO user_courses (user_id, course_id) 
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
(3, 1), (3, 2), (3, 3),
(4, 7), (5, 8), (6, 9), (7, 10),
(9, 1), (9, 2), (9, 3), (9, 4), (9, 5), (9, 6), (9, 7), (9, 8), (9, 9), (9, 10);