-- =============================================
-- SQL Script for Database Initialization
-- Created: 2026
-- =============================================

-- =============================================
-- 1. CREATE USERS/ROLES
-- =============================================

CREATE ROLE article_service_role WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    PASSWORD '{{ .Env.ARTICLE_SERVICE_PASSWORD }}'
    CONNECTION LIMIT 50;

-- =============================================
-- 2. CREATE DATABASES
-- =============================================

-- Main application database
CREATE DATABASE article_service_db
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

\c article_service_db

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Set default privileges
GRANT CONNECT ON DATABASE article_service_db TO article_service_role;
GRANT USAGE ON SCHEMA public TO article_service_role;

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO article_service_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO article_service_role;

-- =============================================
-- 4. CREATE TYPES
-- =============================================

CREATE TYPE article_status AS ENUM (
    'published',
    'draft',
    'archived'
);

-- =============================================
-- 5. CREATE TABLES
-- =============================================

CREATE TABLE metadata (
uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    
title VARCHAR UNIQUE NOT NULL,
status article_status NOT NULL DEFAULT 'draft',
text_id VARCHAR UNIQUE,
video_url VARCHAR UNIQUE,

created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
 
);
-- =============================================
-- 6. FUNCTIONS
-- =============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- =============================================
-- 7. TRIGGERS
-- =============================================
CREATE TRIGGER metadata_updated_at
    BEFORE UPDATE ON metadata
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
