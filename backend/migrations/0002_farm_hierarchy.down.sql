-- Drop tables in reverse order of creation to respect foreign key constraints

DROP TABLE IF EXISTS plot_crops;
DROP TABLE IF EXISTS crops;
DROP TABLE IF EXISTS plot_leases;
DROP TABLE IF EXISTS plots;
DROP TABLE IF EXISTS greenhouses;
DROP TABLE IF EXISTS land_parcels;
DROP TABLE IF EXISTS regions;
