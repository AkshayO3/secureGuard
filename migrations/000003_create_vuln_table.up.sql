CREATE TYPE vulnerability_severity AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE vulnerability_status AS ENUM ('open', 'in_progress', 'mitigated', 'resolved');

CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cve_id VARCHAR(50),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity vulnerability_severity NOT NULL,
    status vulnerability_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);