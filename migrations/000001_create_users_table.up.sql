CREATE TYPE user_role AS ENUM ('admin', 'analyst', 'viewer');

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username VARCHAR(255) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       role user_role NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);