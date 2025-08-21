-- Создаем новую таблицу для камер
CREATE TABLE IF NOT EXISTS cameras (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    rtsp_path_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_plot
        FOREIGN KEY(plot_id)
        REFERENCES plots(id)
        ON DELETE CASCADE
);

-- Добавляем индекс для plot_id для ускорения выборок
CREATE INDEX IF NOT EXISTS idx_cameras_on_plot_id ON cameras(plot_id);

-- Удаляем старое поле camera_url из таблицы plots
ALTER TABLE plots DROP COLUMN IF EXISTS camera_url;
