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
('Компьютерный практикум', 1, 1, 
'## Комптютерный практикум \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=1) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=1) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=1)', 2),

('Информационные системы и технологии', 2, 1, '## Информационные системы и технологии \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=2) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=2) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=2)', 2),

('Информационные системы и технологии', 3, 1, '## Информационные системы и технологии \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=3) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=3) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=3)', 2),

('Философия', 4, 1, '## Философия \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=4) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=4) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=4)', 6),

('Психология', 5, 1, '## Психология \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=5) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=5) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=5)', 6),

('Большевистская железная дорога', 6, 1, '## БЖД \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=6) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=6) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=6)', 6),

('Крутой предмет О6', 2, 2, '## Предмет О6 \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=7) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=7) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=7)', 3),

('Крутой предмет О4', 2, 3, '## Предмет О4 \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=8) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=8) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=8)', 4),

('Крутой предмет И9', 2, 4, '## Предмет И9 \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=9) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=9) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=9)', 5),

('Крутой предмет П1', 2, 5, '## Предмет П1 \[Технологическая карта](http://134.209.230.107:8080/api/courses/infos?id=10) \[Тест 1](http://134.209.230.107:8080/api/courses/tests?id=10) \[Лабораторная работа 1](http://134.209.230.107:8080/api/courses/labs?id=10)', 6);

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
(1, 'Требования КП', '## Заголовок тербований КП'),
(2, 'Требования ИСИТ', '## Заголовок тербований ИСИТ'),
(3, 'Требования ИСИТ2', '## Заголовок тербований ИСИТ2'),
(4, 'Требования Фиолософия', '## Заголовок тербований Философия'),
(5, 'Требования Психология', '## Заголовок тербований Психология'),
(6, 'Требования БЖД', '## Заголовок тербований большевистской железной дороги'),
(7, 'Требования предмета О6', '## Заголовок тербований предмета О6'),
(8, 'Требования предмета О4', '## Заголовок тербований предмета О4'),
(9, 'Требования предмета И9', '## Заголовок тербований предмета И9'),
(10, 'Требования предмета П1', '## Заголовок тербований предмета П1');

INSERT INTO nested_tests 
(course_id, opens, closes, 
tasks_count, topic, location_id, 
attempts, password, time_limit)
VALUES 
(1, '2024-01-01 08:00:00', '2024-02-01 00:00:00',
15, 'Утилита поэтапной компиляции Make', 2, 1, 'Пароль', '00:15:00'),
(2, '2024-02-01 08:00:00', '2024-03-01 00:00:00',
20, 'C# тест по лабораторной работе 1', 1, 3, '', '00:20:00'),
(3, '2024-03-01 08:00:00', '2024-04-01 00:00:00',
25, 'C++ тест по лабораторной работе 1', 2, 1, 'Пароль', '00:25:00'),
(4, '2024-04-01 08:00:00', '2024-05-01 00:00:00',
15, 'Цитаты Джейсона Стетхема', 1, 3, '', '00:15:00'),
(5, '2024-05-01 08:00:00', '2024-06-01 00:00:00',
20, 'Контроль просмотра сериала солдаты', 2, 1, '', '00:20:00'),
(6, '2024-06-01 08:00:00', '2024-07-01 00:00:00',
25, 'Похвала советского союза', 1, 3, 'Пароль', '00:25:00'),
(7, '2024-07-01 08:00:00', '2024-08-01 00:00:00',
15, 'Тест 1 предмета О6', 2, 1, '', '00:15:00'),
(8, '2024-08-01 08:00:00', '2024-09-01 00:00:00',
20, 'Тест 1 предмета О4', 1, 3, '', '00:20:00'),
(9, '2024-09-01 08:00:00', '2024-10-01 00:00:00',
25, 'Тест 1 предмета И9', 1, 3, '', '00:25:00'),
(10, '2024-10-01 08:00:00', '2024-11-01 00:00:00',
15, 'Тест 1 предмета П1', 2, 1, 'Пароль', '00:15:00');

INSERT INTO nested_labs 
(course_id, opens, closes, topic, requirements, 
example, location_id, attempts)
VALUES
(1, '2024-01-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1: Make', 
'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', '', 2, 1),

(2, '2024-02-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1: C#', 
'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 1, 3),

(3, '2024-03-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1: C++', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 2, 1),

(4, '2024-04-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1: Философия ворониных', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 1, 3),

(5, '2024-05-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1: Повторение опасных трюков из сериалов', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 'https://www.google.com', 2, 1),

(6, '2024-06-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1: Критика запада', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 1, 3),

(7, '2024-07-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1 предмета О6 ', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 2, 1),

(8, '2024-08-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1 предмета О4', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 1, 3),

(9, '2024-09-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1 предмета И9', 'https://docs.google.com/document/d/1r6a0xbuxaqbIAG25Bcu9iBH6t9eCQr0zqJTEEC5VIsg/edit#heading=h.ibwrrmm3ajwe', 
'https://www.google.com', 2, 1),

(10, '2024-10-01 08:00:00', '2024-05-30 00:00:00',
'Лабораторная работа 1 предмета П1', '', '', 1, 3);
