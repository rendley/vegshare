-- Добавляем обратно поле camera_url в таблицу plots
-- Мы не можем знать исходное значение, поэтому добавляем с default ''
ALTER TABLE plots ADD COLUMN camera_url VARCHAR(255) NOT NULL DEFAULT '';

-- Удаляем таблицу cameras
DROP TABLE IF EXISTS cameras;
