CREATE TYPE incident_category AS ENUM ('malware', 'phishing', 'unauthorized_access');
CREATE TYPE incident_severity AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE incident_status AS ENUM ('detected', 'investigating', 'contained', 'resolved');

CREATE TABLE incidents (
                           id UUID PRIMARY KEY,
                           title VARCHAR(255) NOT NULL,
                           description TEXT,
                           category incident_category NOT NULL,
                           severity incident_severity NOT NULL,
                           status incident_status NOT NULL,
                           reported_by UUID NOT NULL REFERENCES users(id),
                           assigned_to UUID REFERENCES users(id),
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);