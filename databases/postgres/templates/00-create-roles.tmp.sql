-- =============================================
-- SQL Script for Database Initialization
-- Created: 2026
-- =============================================

-- =============================================
-- 1. FUNCTIONS
-- =============================================


-- Выдача привилегий роли

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


