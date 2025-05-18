-- +goose Up

-- Расширение
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Функция set_updated_at
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


-- Таблица account
CREATE TABLE account (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    surname VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(60),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);
CREATE TRIGGER trigger_set_updated_at_on_account
BEFORE UPDATE ON account
FOR EACH ROW EXECUTE FUNCTION set_updated_at();


-- Начальные данные
INSERT INTO account (name, surname, email, password_hash)
VALUES 
    ('админ', 'админ', 'test@yandex.ru', '$2a$04$hBpYvS1On29wYe6yPLdBL.EoWH0WoD7JGgQaIxV4mA.TTXEsAXfdO');



-- +goose Down

-- Удаление account
DROP TRIGGER IF EXISTS trigger_set_updated_at_on_account ON account;
DROP TABLE IF EXISTS account;

-- Удаление функции
DROP FUNCTION IF EXISTS set_updated_at;

-- Удаление расширения
DROP EXTENSION IF EXISTS "uuid-ossp";
