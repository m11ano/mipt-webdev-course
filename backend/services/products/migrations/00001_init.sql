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

-- Таблица файлов
CREATE TABLE file (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_target             SMALLINT NOT NULL,
    assigned_to_target      BOOLEAN NOT NULL DEFAULT FALSE,
    storage_file_key        VARCHAR(100) NULL,
    uploaded_to_storage     BOOLEAN NOT NULL DEFAULT FALSE,
    to_delete_from_storage  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NULL,
    deleted_at              TIMESTAMPTZ NULL
);
CREATE INDEX idx_file_storage_file_key ON file(storage_file_key);
CREATE TRIGGER trigger_set_updated_at_on_image
BEFORE UPDATE ON file
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Таблица товаров
CREATE TABLE product (
    id                      BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    is_published            BOOLEAN NOT NULL DEFAULT FALSE,
    name                    VARCHAR(150) NOT NULL,
    full_description        TEXT NOT NULL,
    price                   NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    stock_available         INTEGER NOT NULL DEFAULT 0 CHECK (stock_available >= 0),
    image_preview_file_id   UUID REFERENCES file(id) ON DELETE SET NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NULL,
    deleted_at              TIMESTAMPTZ NULL
);
CREATE INDEX idx_product_created_at ON product(created_at);
CREATE TRIGGER trigger_set_updated_at_on_product
BEFORE UPDATE ON product
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Таблица product_order_block
CREATE TABLE product_order_block (
    product_id BIGINT NOT NULL REFERENCES product(id) ON DELETE RESTRICT,
    order_id   BIGINT NOT NULL,
    quantity   INTEGER NOT NULL CHECK (quantity >= 1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (product_id, order_id)
);
CREATE INDEX idx_product_order_block_order_id ON product_order_block(order_id);

-- Таблица product_slider_image
CREATE TABLE product_slider_image (
    product_id  BIGINT NOT NULL REFERENCES product(id) ON DELETE RESTRICT,
    file_id     UUID NOT NULL REFERENCES file(id) ON DELETE CASCADE,
    sort        INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (product_id, file_id)
);

-- +goose Down

-- Удаление product_slider_image
DROP TABLE IF EXISTS product_slider_image;

-- Удаление product_order_block
DROP INDEX IF EXISTS idx_product_order_block_order_id;
DROP TABLE IF EXISTS product_order_block;

-- Удаление товаров
DROP INDEX IF EXISTS idx_product_created_at;
DROP TRIGGER IF EXISTS trigger_set_updated_at_on_product ON product;
DROP TABLE IF EXISTS product;

-- Удаление file
DROP INDEX IF EXISTS idx_file_storage_file_key;
DROP TRIGGER IF EXISTS trigger_set_updated_at_on_file ON file;
DROP TABLE IF EXISTS file;


-- Удаление функций
DROP FUNCTION IF EXISTS set_updated_at;

-- Удаление расширения
DROP EXTENSION IF EXISTS "uuid-ossp";
