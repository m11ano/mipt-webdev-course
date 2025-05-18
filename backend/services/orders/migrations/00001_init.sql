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

-- Таблица заказов
CREATE TABLE "order" (
    id                      BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    status                  INTEGER NOT NULL DEFAULT 0,
    order_sum               NUMERIC(12, 2) NOT NULL CHECK (order_sum >= 0),
    secret_key              UUID NOT NULL DEFAULT uuid_generate_v4(),
    client_name             VARCHAR(150) NOT NULL,
    client_surname          VARCHAR(150) NOT NULL,
    client_email            VARCHAR(150) NOT NULL,
    client_phone            VARCHAR(20) NOT NULL,
    delivery_address        VARCHAR(150) NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NULL,
    deleted_at              TIMESTAMPTZ NULL
);
CREATE TRIGGER trigger_set_updated_at_on_order
BEFORE UPDATE ON "order"
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Таблица order_product
CREATE TABLE order_product (
    order_id        BIGINT NOT NULL REFERENCES "order"(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL,
    price           NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    quantity        INTEGER NOT NULL CHECK (quantity >= 1),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (order_id, product_id)
);
CREATE INDEX idx_product_order_block_product_id ON order_product(product_id);

-- +goose Down

-- Удаление order_product
DROP INDEX IF EXISTS idx_product_order_block_product_id;
DROP TABLE IF EXISTS order_product;

-- Удаление заказов
DROP TRIGGER IF EXISTS trigger_set_updated_at_on_order ON "order";
DROP TABLE IF EXISTS "order";


-- Удаление функций
DROP FUNCTION IF EXISTS set_updated_at;

-- Удаление расширения
DROP EXTENSION IF EXISTS "uuid-ossp";
