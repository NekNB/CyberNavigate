-- =============================================
-- SQL Script for Database Initialization
-- Created: 2026
-- =============================================

-- =============================================
-- 1. CREATE USERS/ROLES
-- =============================================

CREATE ROLE user_service_role WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    PASSWORD '{{ .Env.USER_SERVICE_PASSWORD }}'
    CONNECTION LIMIT 50;

-- =============================================
-- 2. CREATE DATABASES
-- =============================================

-- Main application database
CREATE DATABASE user_service_db
    WITH
    OWNER = {{ .Env.POSTGRES_USER }}
    ENCODING 'UTF8'
    LC_COLLATE 'en_US.UTF-8'
    LC_CTYPE 'en_US.UTF-8'
    TEMPLATE template0
    CONNECTION LIMIT -1;



-- =============================================
-- 3. CONNECT TO MAIN DATABASE AND SETUP SCHEMA
-- =============================================

\c user_service_db

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Set default privileges
GRANT CONNECT ON DATABASE user_service_db TO user_service_role;
GRANT USAGE ON SCHEMA public TO user_service_role;

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO user_service_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO user_service_role;



-- =============================================
-- 4. CREATE TABLES
-- =============================================

-- Users table
CREATE TABLE users (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    salt VARCHAR(256) NOT NULL,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(uuid),
    refresh_token VARCHAR(16) NOT NULL UNIQUE DEFAULT encode(gen_random_bytes(8), 'hex'),
    revoked BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
 

-- =============================================
-- 5. FUNCTIONS
-- =============================================


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


-- =============================================
-- 6. TRIGGERS
-- =============================================

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
