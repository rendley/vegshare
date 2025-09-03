-- 1. Rename the main table
ALTER TABLE greenhouses RENAME TO structures;

-- 2. Rename the foreign key column in the plots table
ALTER TABLE plots RENAME COLUMN greenhouse_id TO structure_id;

-- 3. Rename the foreign key constraint itself
ALTER TABLE plots RENAME CONSTRAINT plots_greenhouse_id_fkey TO plots_structure_id_fkey;