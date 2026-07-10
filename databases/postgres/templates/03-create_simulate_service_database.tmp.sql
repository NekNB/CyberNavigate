

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

\c simulate_service_db

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

CREATE TABLE difficulty (
    level VARCHAR(10) PRIMARY KEY
);

CREATE TABLE action_type (
    name VARCHAR(10) PRIMARY KEY
);


CREATE TABLE answers (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    text VARCHAR(50) NOT NULL,
    add_trust INTEGER NOT NULL DEFAULT 0,
    error VARCHAR(256)
);



CREATE TABLE steps (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    previous_step UUID REFERENCES steps(uuid),
    previous_answer UUID REFERENCES answers(uuid),
    min_trust INT,
    max_trust INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_not_both_not_null CHECK (previous_step IS NULL OR previous_answer IS NULL),

    CONSTRAINT chk_max_greater_min CHECK (max_trust > min_trust)
);

CREATE TABLE scenarios (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    title VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200) NOT NULL,
    difficulty VARCHAR(10) NOT NULL REFERENCES difficulty(level),
    first_step UUID UNIQUE REFERENCES steps(uuid),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE article_ids_to_scenario_id (
    article_id UUID PRIMARY KEY,
    senario_id UUID NOT NULL REFERENCES scenarios(uuid)
);
CREATE TABLE actions (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    step_id UUID NOT NULL REFERENCES steps(uuid),
    type VARCHAR(10) NOT NULL REFERENCES action_type(name),
    message_id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    delay INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE files (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    filename VARCHAR(30) NOT NULL UNIQUE,
    message_id UUID NOT NULL REFERENCES messages(uuid), 
    is_safe BOOLEAN NOT NULL,
    size INTEGER NOT NULL
    error VARCHAR(256)
);



CREATE TABLE messages (
    uuid UUID PRIMARY KEY REFERENCES actions(object_id),
    sender_id UUID NOT NULL,
    sender_name VARCHAR(50) NOT NULL, 
    text VARCHAR(1024),
);

CREATE TABLE errors_to_user (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    session_id UUID REFERENCES session(uuid)
    error VARCHAR(256) NOT NULL,
)





CREATE TABLE sessions (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    current_trust INTEGER CHECK (current_trust BETWEEN -100 AND 100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE trust (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(uuid),
    noted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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

CREATE TRIGGER update_scenarios_updated_at
    BEFORE UPDATE ON scenarios
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_steps_updated_at
    BEFORE UPDATE ON steps
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- 7. INSERT DATA
-- =============================================
INSERT INTO difficulty (level) VALUES 
    ('easy'),
    ('middle'),
    ('hard');


INSERT INTO action_type (name) VALUES 
    ('message'),
    ('sms');