CREATE DATABASE LinkDB;
ALTER ROLE admin SET client_encoding TO 'utf8';
GRANT ALL PRIVILEGES ON DATABASE LinkDB TO admin;

-- Таблица для отображения сокращенных ссылок на оригинальные ссылки
CREATE TABLE IF NOT EXISTS link_mapping
(
    short_link VARCHAR(10) PRIMARY KEY, -- Сокращенная ссылка
    original_link TEXT NOT NULL -- Оригинальная ссылка
);