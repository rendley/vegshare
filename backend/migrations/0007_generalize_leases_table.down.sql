-- Откатываем изменения в обратном порядке
ALTER TABLE leases RENAME COLUMN unit_id TO plot_id;
ALTER TABLE leases DROP COLUMN unit_type;
ALTER TABLE leases RENAME TO plot_leases;
