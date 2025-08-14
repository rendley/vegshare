-- Начинаем транзакцию.
BEGIN;

-- Удаляем таблицы в порядке, обратном их созданию,
-- чтобы не нарушать ограничения внешних ключей.
DROP TABLE IF EXISTS plot_crops;
DROP TABLE IF EXISTS crops;
DROP TABLE IF EXISTS plots;
DROP TABLE IF EXISTS farms;

-- Завершаем транзакцию.
COMMIT;
