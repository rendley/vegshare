DROP TABLE operation_log;

-- Junction Table for what crop is planted on which plot
CREATE TABLE plot_crops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plot_id UUID NOT NULL, -- REFERENCES plots(id) ON DELETE CASCADE,
    crop_id UUID NOT NULL, -- REFERENCES crops(id) ON DELETE RESTRICT,
    lease_id UUID NOT NULL, -- REFERENCES plot_leases(id) ON DELETE CASCADE,
    planted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'growing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON plot_crops (plot_id);
CREATE INDEX ON plot_crops (crop_id);
CREATE INDEX ON plot_crops (lease_id);
