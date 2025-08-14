-- Начинаем транзакцию.
BEGIN;

-- Таблица для хранения информации о тепличных комплексах (фермах).
CREATE TABLE farms (
    -- Уникальный идентификатор фермы.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Название фермы, например, "Теплица в Подмосковье".
    name VARCHAR(255) NOT NULL,
    -- Местоположение или адрес фермы.
    location TEXT
);

-- Таблица для хранения информации об участках земли на фермах.
CREATE TABLE plots (
    -- Уникальный идентификатор участка.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Внешний ключ, связывающий участок с конкретной фермой.
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    -- Внешний ключ, связывающий участок с пользователем, который его арендует.
    -- Может быть NULL, если участок свободен.
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    -- Размер участка, например, в квадратных метрах.
    size REAL NOT NULL,
    -- Дата и время создания записи.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Справочник доступных для посадки культур.
CREATE TABLE crops (
    -- Уникальный идентификатор культуры.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Название культуры, например, "Томат Черри" или "Огурец Родничок".
    name VARCHAR(255) UNIQUE NOT NULL,
    -- Описание особенностей выращивания.
    description TEXT,
    -- Среднее время роста в днях.
    growing_time_days INTEGER
);

-- Таблица для отслеживания того, какая культура посажена на каком участке.
CREATE TABLE plot_crops (
    -- Уникальный идентификатор посадки.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Внешний ключ, связывающий с участком.
    plot_id UUID NOT NULL REFERENCES plots(id) ON DELETE CASCADE,
    -- Внешний ключ, связывающий с культурой из справочника.
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE RESTRICT,
    -- Статус посадки: 'planted', 'growing', 'ready_for_harvest', 'harvested'.
    status VARCHAR(50) NOT NULL DEFAULT 'planted',
    -- Дата посадки.
    planted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создаем индексы для ускорения поиска по внешним ключам.
CREATE INDEX idx_plots_farm_id ON plots(farm_id);
CREATE INDEX idx_plots_owner_id ON plots(owner_id);
CREATE INDEX idx_plot_crops_plot_id ON plot_crops(plot_id);
CREATE INDEX idx_plot_crops_crop_id ON plot_crops(crop_id);

-- Завершаем транзакцию.
COMMIT;
