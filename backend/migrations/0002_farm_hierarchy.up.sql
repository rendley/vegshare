-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table for Regions
CREATE TABLE regions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table for Land Parcels (previously Farms)
CREATE TABLE land_parcels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table for Greenhouses
CREATE TABLE greenhouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    land_parcel_id UUID NOT NULL REFERENCES land_parcels(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100), -- e.g., 'hydroponic', 'soil'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table for Plots (the rentable units)
CREATE TABLE plots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    greenhouse_id UUID NOT NULL REFERENCES greenhouses(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    size VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'available', -- e.g., 'available', 'rented', 'maintenance'
    camera_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table for Plot Leases (connecting users to plots)
CREATE TABLE plot_leases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plot_id UUID NOT NULL REFERENCES plots(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- e.g., 'active', 'expired', 'cancelled'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(plot_id, user_id, status) -- A user can only have one active lease per plot
);

-- Table for Crops (the catalog of plants)
CREATE TABLE crops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    -- Time in days
    planting_time INT,
    harvest_time INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Junction Table for what crop is planted on which plot
CREATE TABLE plot_crops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plot_id UUID NOT NULL REFERENCES plots(id) ON DELETE CASCADE,
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE RESTRICT,
    lease_id UUID NOT NULL REFERENCES plot_leases(id) ON DELETE CASCADE, -- Link to the specific lease
    planted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'growing', -- e.g., 'growing', 'harvested', 'failed'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add indexes for foreign keys to improve query performance
CREATE INDEX ON land_parcels (region_id);
CREATE INDEX ON greenhouses (land_parcel_id);
CREATE INDEX ON plots (greenhouse_id);
CREATE INDEX ON plot_leases (plot_id);
CREATE INDEX ON plot_leases (user_id);
CREATE INDEX ON plot_crops (plot_id);
CREATE INDEX ON plot_crops (crop_id);
CREATE INDEX ON plot_crops (lease_id);
