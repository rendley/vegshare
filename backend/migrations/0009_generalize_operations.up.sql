-- 1. Переименовываем таблицу
ALTER TABLE plot_crops RENAME TO operation_log;

-- 2. Удаляем старые, ненужные колонки
ALTER TABLE operation_log DROP COLUMN crop_id;
ALTER TABLE operation_log DROP COLUMN lease_id;

-- 3. Переименовываем и добавляем новые колонки
ALTER TABLE operation_log RENAME COLUMN plot_id TO unit_id;
ALTER TABLE operation_log RENAME COLUMN planted_at TO executed_at;
ALTER TABLE operation_log ADD COLUMN unit_type VARCHAR(50) NOT NULL;
ALTER TABLE operation_log ADD COLUMN user_id UUID NOT NULL;
ALTER TABLE operation_log ADD COLUMN action_type VARCHAR(50) NOT NULL;
ALTER TABLE operation_log ADD COLUMN parameters JSONB;

-- 4. Добавляем внешние ключи и индексы
ALTER TABLE operation_log ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX ON operation_log (user_id);
CREATE INDEX ON operation_log (unit_id, unit_type);
