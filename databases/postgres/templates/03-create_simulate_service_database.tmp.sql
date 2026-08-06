

-- =============================================
-- SQL Script for Database Initialization
-- Created: 2026
-- =============================================

-- =============================================
-- 1. CREATE USERS/ROLES
-- =============================================

CREATE ROLE simulator_service_role WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    PASSWORD '{{ .Env.SIMULATOR_SERVICE_PASSWORD }}'
    CONNECTION LIMIT 50;

-- =============================================
-- 2. CREATE DATABASES
-- =============================================

-- Main application database
CREATE DATABASE simulator_service_db
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

\c simulator_service_db

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Set default privileges
GRANT CONNECT ON DATABASE simulator_service_db TO simulator_service_role;
GRANT USAGE ON SCHEMA public TO simulator_service_role;

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO simulator_service_role;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO simulator_service_role;



-- =============================================
-- 4. CREATE TABLES
-- =============================================

CREATE TABLE difficulty (
    level VARCHAR(10) PRIMARY KEY
);

CREATE TABLE action_type (
    name VARCHAR(10) PRIMARY KEY
);


CREATE TABLE messages (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    sender_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    sender_name VARCHAR(50) NOT NULL, 
    text VARCHAR(1024)
);

CREATE TABLE files (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    filename VARCHAR(30) NOT NULL, 
    is_safe BOOLEAN NOT NULL,
    message_id UUID  REFERENCES messages(uuid) ON DELETE CASCADE,
    size INTEGER NOT NULL,
    error VARCHAR(256)
);
CREATE TABLE answers (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    text VARCHAR(50) NOT NULL,
    message_id UUID REFERENCES messages(uuid) ON DELETE CASCADE,
    add_trust INTEGER NOT NULL DEFAULT 0,
    error VARCHAR(256)
);



CREATE TABLE steps (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   
    min_trust INT NOT NULL DEFAULT -100,
    max_trust INT NOT NULL DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_not_both_not_null CHECK (previous_step IS NULL OR previous_answer IS NULL),

    CONSTRAINT chk_max_greater_min CHECK (max_trust > min_trust),
    CONSTRAINT chk_max_normal CHECK (max_trust BETWEEN -100 AND 100),
    CONSTRAINT chk_min_normal CHECK (min_trust BETWEEN -100 AND 100)
);

CREATE TABLE step_to_steps (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    step_uuid UUID NOT NULL REFERENCES steps(uuid) ON DELETE CASCADE,
    previous_step_uuid UUID NOT NULL REFERENCES steps(uuid) ON DELETE CASCADE
)
CREATE TABLE step_to_answers (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    step_uuid UUID NOT NULL REFERENCES steps(uuid) ON DELETE CASCADE,
    previous_answer_uuid UUID NOT NULL REFERENCES answers(uuid) ON DELETE CASCADE
)

CREATE TABLE scenarios (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    title VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200) NOT NULL,
    difficulty VARCHAR(10) NOT NULL REFERENCES difficulty(level) ON DELETE CASCADE,
    first_step UUID UNIQUE REFERENCES steps(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE steps ADD COLUMN scenario_id UUID NOT NULL REFERENCES scenarios(uuid) ON DELETE CASCADE;

CREATE TABLE article_ids_to_scenario_id (
    article_id UUID PRIMARY KEY,
    scenario_id UUID NOT NULL REFERENCES scenarios(uuid ) ON DELETE CASCADE
);


CREATE TABLE actions (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    step_id UUID REFERENCES steps(uuid) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL REFERENCES action_type(name) ON DELETE CASCADE,
    message_id UUID NOT NULL REFERENCES messages(uuid) ON DELETE CASCADE,
    delay INTEGER NOT NULL DEFAULT 0
);





CREATE TABLE sessions (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL,
    current_trust INTEGER CHECK (current_trust BETWEEN -100 AND 100) DEFAULT 0,
    current_step UUID REFERENCES steps(uuid) ON DELETE CASCADE,
    is_finished BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE errors_to_session (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    session_id UUID REFERENCES sessions(uuid) ON DELETE CASCADE,
    error VARCHAR(256) NOT NULL
);
CREATE TABLE trusts (
    uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    trust INTEGER NOT NULL CHECK (trust BETWEEN -100 AND 100),
    session_id UUID NOT NULL REFERENCES sessions(uuid) ON DELETE CASCADE,
    noted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
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

-- =============================================
-- 8. INDEX
-- =============================================

-- Гарантируем, что у пользователя может быть только одна активная сессия
CREATE UNIQUE INDEX idx_sessions_active_user 
ON sessions (user_id) 
WHERE finished_at IS NULL;