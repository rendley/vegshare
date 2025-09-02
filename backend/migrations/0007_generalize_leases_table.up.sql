-- Переименовываем таблицу
ALTER TABLE plot_leases RENAME TO leases;

-- Переименовываем колонку
ALTER TABLE leases RENAME COLUMN plot_id TO unit_id;

-- Добавляем новую колонку для типа юнита
ALTER TABLE leases ADD COLUMN unit_type VARCHAR(50);

-- Заполняем значение 'plot' для всех существующих записей
UPDATE leases SET unit_type = 'plot';

-- Устанавливаем ограничение NOT NULL, так как теперь все поля заполнены
ALTER TABLE leases ALTER COLUMN unit_type SET NOT NULL;
