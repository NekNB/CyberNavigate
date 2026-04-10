-- =============================================
-- SQL Script for Database Initialization
-- Created: 2026
-- =============================================

-- =============================================
-- 1. CREATE USERS/ROLES
-- =============================================

CREATE ROLE simulate_service_role WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    PASSWORD '{{ .Env.SIMULATE_SERVICE_PASSWORD }}'
    CONNECTION LIMIT 50;

-- =============================================
-- 2. CREATE DATABASES
-- =============================================

-- Main application database
CREATE DATABASE simulate_service_db
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
GRANT CONNECT ON DATABASE user_service_db TO simulate_service_role;
GRANT USAGE ON SCHEMA public TO simulate_service_role;

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO simulate_service_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO simulate_service_role;



-- =============================================
-- 4. CREATE TABLES
-- =============================================


 

-- =============================================
-- 5. FUNCTIONS
-- =============================================





-- =============================================
-- 6. TRIGGERS
-- =============================================

-- CREATE TRIGGER update_users_updated_at
--     BEFORE UPDATE ON users
--     FOR EACH ROW
--     EXECUTE FUNCTION update_updated_at_column();
