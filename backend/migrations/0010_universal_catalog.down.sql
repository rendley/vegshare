-- Удаляем новую таблицу каталога
DROP TABLE IF EXISTS catalog_items;

-- Воссоздаем старую таблицу crops для возможности отката
CREATE TABLE crops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    planting_time INT,
    harvest_time INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
