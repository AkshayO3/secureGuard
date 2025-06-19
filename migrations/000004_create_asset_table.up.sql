CREATE TYPE asset_type AS ENUM ('server', 'workstation', 'network_device', 'application');
CREATE TYPE asset_status AS ENUM ('active', 'inactive', 'decommissioned');

CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type asset_type NOT NULL,
    ip_address VARCHAR(45),
    os VARCHAR(255),
    status asset_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);