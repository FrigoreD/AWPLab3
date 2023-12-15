ИДБ-20-11 Агапов А. А.

Написание приложения опросника с заранее подготовленными вопросами и ответами

Развёртывание:
1. Поднял локально MySQL, сделал подключение к базе в main().  Добавил таблицу и записи в неё: 
CREATE TABLE IF NOT EXISTS questions (
id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
question VARCHAR(255) NOT NULL,
answer1 VARCHAR(100) NOT NULL,
answer2 VARCHAR(100) NOT NULL,
correct_answer TINYINT(1) UNSIGNED NOT NULL
);

INSERT INTO questions (question, answer1, answer2, correct_answer) VALUES
('Вопрос 1', 'Ответ 1', 'Другой ответ', 1),
('Вопрос 2', 'Ответ A', 'Ответ B', 0);

2. Увидел что можно сделать через fyne либу. Поднимал интерфейс по этой документации: https://developer.fyne.io/started/. Тут нужен gcc для компиляции интерфейса на C. Немного повозился.
