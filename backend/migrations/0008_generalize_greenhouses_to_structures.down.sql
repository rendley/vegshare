-- Revert in reverse order

-- 3. Rename the foreign key constraint back
ALTER TABLE plots RENAME CONSTRAINT plots_structure_id_fkey TO plots_greenhouse_id_fkey;

-- 2. Rename the foreign key column back
ALTER TABLE plots RENAME COLUMN structure_id TO greenhouse_id;

-- 1. Rename the main table back
ALTER TABLE structures RENAME TO greenhouses;