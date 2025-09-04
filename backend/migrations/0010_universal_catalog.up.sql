-- Создаем новую универсальную таблицу для каталога
CREATE TABLE catalog_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_type VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    attributes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(item_type, name)
);

-- Добавляем индексы для ускорения выборок
CREATE INDEX idx_catalog_items_item_type ON catalog_items(item_type);
CREATE INDEX idx_catalog_items_name ON catalog_items(name);

-- Удаляем старую таблицу, так как она больше не нужна
DROP TABLE IF EXISTS crops;
