-- ============================================================================
-- COMPLETE USER MANAGEMENT & AUTHENTICATION SCHEMA
-- PostgreSQL Database Schema
-- Version: 2.0
-- Last Updated: January 14, 2026
-- ============================================================================

-- ============================================================================
-- EXTENSIONS
-- ============================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy text search
CREATE EXTENSION IF NOT EXISTS "btree_gin"; -- For GIN indexes

-- ============================================================================
-- CUSTOM TYPES
-- ============================================================================
CREATE TYPE user_role AS ENUM ('admin', 'leader', 'user');
CREATE TYPE account_status_type AS ENUM ('pending', 'active', 'suspended', 'banned', 'closed');
CREATE TYPE severity_type AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE social_provider AS ENUM ('google', 'facebook', 'apple');
CREATE TYPE platform_type AS ENUM ('ios', 'android', 'web');
CREATE TYPE theme_type AS ENUM ('light', 'dark', 'auto');
CREATE TYPE email_status AS ENUM ('pending', 'sending', 'sent', 'failed');

-- ============================================================================
-- USERS TABLE (Core Entity)
-- ============================================================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    role user_role NOT NULL,
    
    -- Account status
    account_status account_status_type NOT NULL DEFAULT 'pending',
    status_changed_at TIMESTAMP,
    status_changed_by INT REFERENCES users(id) ON DELETE SET NULL,
    status_reason TEXT,
    
    -- Profile
    phone VARCHAR(20),
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    
    -- Emergency info
    blood_type VARCHAR(10),
    allergies TEXT,
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relationship VARCHAR(50),
    
    -- Verification status
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP,
    phone_verified BOOLEAN DEFAULT FALSE,
    phone_verified_at TIMESTAMP,
    
    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id) ON DELETE SET NULL,
    last_active_at TIMESTAMP
);

-- Indexes for users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_account_status ON users(account_status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_email_verified ON users(email_verified);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_last_active_at ON users(last_active_at);

-- Full-text search indexes
CREATE INDEX idx_users_name_gin ON users USING gin(to_tsvector('english', name));
CREATE INDEX idx_users_email_gin ON users USING gin(to_tsvector('english', email));

-- ============================================================================
-- USER CREDENTIALS (Separated for Security)
-- ============================================================================
CREATE TABLE user_credentials (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash VARCHAR(255),
    
    -- Two-factor authentication
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_secret VARCHAR(255), -- Encrypted
    two_factor_enabled_at TIMESTAMP,
    
    -- Password policy
    must_change_password BOOLEAN DEFAULT FALSE,
    password_expires_at TIMESTAMP,
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP,
    deleted_by INT REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_credentials_user_id ON user_credentials(user_id);
CREATE INDEX idx_credentials_must_change ON user_credentials(must_change_password);

-- ============================================================================
-- USER SECURITY INFO (Separated for Performance)
-- ============================================================================
CREATE TABLE user_security_info (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Login attempts
    failed_login_attempts INT DEFAULT 0,
    last_failed_login_at TIMESTAMP,
    locked_until TIMESTAMP,
    
    -- Password management
    last_password_change TIMESTAMP,
    password_changed_count INT DEFAULT 0,
    
    -- Activity tracking
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(50),
    last_login_location VARCHAR(255),
    
    -- Security score (0-100)
    security_score INT DEFAULT 50,
    security_score_updated_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_security_info_user_id ON user_security_info(user_id);
CREATE INDEX idx_security_info_locked_until ON user_security_info(locked_until);
CREATE INDEX idx_security_info_security_score ON user_security_info(security_score);

-- ============================================================================
-- PASSWORD HISTORY
-- ============================================================================
CREATE TABLE password_history (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_history_user_id ON password_history(user_id);
CREATE INDEX idx_password_history_created_at ON password_history(created_at);

-- ============================================================================
-- PASSWORD RESET TOKENS
-- ============================================================================
CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    
    -- Security audit
    ip_address VARCHAR(50),
    user_agent TEXT,
    used_by_session_id INT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_reset_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_token ON password_reset_tokens(token);
CREATE INDEX idx_password_reset_expires_at ON password_reset_tokens(expires_at);
CREATE INDEX idx_password_reset_used_at ON password_reset_tokens(used_at);

-- ============================================================================
-- EMAIL VERIFICATION TOKENS
-- ============================================================================
CREATE TABLE email_verification_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL, -- Support for email change verification
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    
    -- Security audit
    ip_address VARCHAR(50),
    user_agent TEXT,
    used_by_session_id INT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_token_user_id ON email_verification_tokens(user_id);
CREATE INDEX idx_email_token_token ON email_verification_tokens(token);
CREATE INDEX idx_email_token_expires_at ON email_verification_tokens(expires_at);
CREATE INDEX idx_email_token_email ON email_verification_tokens(email);

-- ============================================================================
-- PHONE VERIFICATION
-- ============================================================================
CREATE TABLE phone_verification_codes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone VARCHAR(20) NOT NULL,
    code VARCHAR(10) NOT NULL, -- Hashed
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    attempts INT DEFAULT 0,
    max_attempts INT DEFAULT 3,
    
    -- Security
    ip_address VARCHAR(50),
    user_agent TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_phone_verification_user_id ON phone_verification_codes(user_id);
CREATE INDEX idx_phone_verification_phone ON phone_verification_codes(phone);
CREATE INDEX idx_phone_verification_expires_at ON phone_verification_codes(expires_at);

-- ============================================================================
-- SOCIAL AUTHENTICATION
-- ============================================================================
CREATE TABLE user_social_auth (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider social_provider NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    
    -- OAuth tokens (should be encrypted in application)
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    
    -- Cached profile
    provider_email VARCHAR(255),
    provider_name VARCHAR(100),
    provider_avatar_url TEXT,
    provider_data JSONB, -- Additional provider-specific data
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_social_auth_user_id ON user_social_auth(user_id);
CREATE INDEX idx_social_auth_provider ON user_social_auth(provider);
CREATE INDEX idx_social_auth_provider_user_id ON user_social_auth(provider_user_id);

-- ============================================================================
-- TWO-FACTOR BACKUP CODES
-- ============================================================================
CREATE TABLE two_factor_backup_codes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL, -- Hashed
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_2fa_backup_user_id ON two_factor_backup_codes(user_id);
CREATE INDEX idx_2fa_backup_used_at ON two_factor_backup_codes(used_at);

-- ============================================================================
-- USER SESSIONS
-- ============================================================================
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    
    -- Device info
    device_id VARCHAR(255),
    device_name VARCHAR(100),
    platform platform_type,
    app_version VARCHAR(50),
    os_version VARCHAR(50),
    
    -- Security
    ip_address VARCHAR(50),
    user_agent TEXT,
    location VARCHAR(255),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Lifecycle
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP,
    revoked_by INT REFERENCES users(id) ON DELETE SET NULL,
    revoke_reason TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_session_token ON user_sessions(session_token);
CREATE INDEX idx_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_sessions_revoked_at ON user_sessions(revoked_at);
CREATE INDEX idx_sessions_device_id ON user_sessions(device_id);
CREATE INDEX idx_sessions_last_used_at ON user_sessions(last_used_at);

-- ============================================================================
-- PUSH TOKENS (Separated from User - One user, multiple devices)
-- ============================================================================
CREATE TABLE user_push_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    platform platform_type NOT NULL,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP,
    deactivation_reason TEXT,
    
    -- Metadata
    app_version VARCHAR(50),
    os_version VARCHAR(50),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, device_id, platform)
);

CREATE INDEX idx_push_tokens_user_id ON user_push_tokens(user_id);
CREATE INDEX idx_push_tokens_device_id ON user_push_tokens(device_id);
CREATE INDEX idx_push_tokens_is_active ON user_push_tokens(is_active);
CREATE INDEX idx_push_tokens_platform ON user_push_tokens(platform);

-- ============================================================================
-- USER PREFERENCES
-- ============================================================================
CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Localization
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    date_format VARCHAR(20) DEFAULT 'YYYY-MM-DD',
    time_format VARCHAR(10) DEFAULT '24h', -- '12h' or '24h'
    
    -- Notifications
    notifications_enabled BOOLEAN DEFAULT TRUE,
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT TRUE,
    sms_notifications BOOLEAN DEFAULT FALSE,
    
    -- Notification preferences by type
    notify_appointments BOOLEAN DEFAULT TRUE,
    notify_messages BOOLEAN DEFAULT TRUE,
    notify_reminders BOOLEAN DEFAULT TRUE,
    notify_promotions BOOLEAN DEFAULT FALSE,
    notify_system_updates BOOLEAN DEFAULT TRUE,
    
    -- Display preferences
    theme theme_type DEFAULT 'light',
    
    -- Privacy
    profile_visibility VARCHAR(20) DEFAULT 'private', -- 'public', 'private', 'contacts'
    show_online_status BOOLEAN DEFAULT TRUE,
    
    -- Other preferences (stored as JSON for flexibility)
    custom_preferences JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_preferences_user_id ON user_preferences(user_id);

-- ============================================================================
-- USER CONSENTS (GDPR Compliance)
-- ============================================================================
CREATE TABLE user_consents (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL, -- 'terms', 'privacy', 'marketing', 'data_processing'
    consent_version VARCHAR(20) NOT NULL,
    granted BOOLEAN NOT NULL,
    
    -- Audit trail
    ip_address VARCHAR(50),
    user_agent TEXT,
    location VARCHAR(255),
    
    granted_at TIMESTAMP,
    revoked_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_consents_user_id ON user_consents(user_id);
CREATE INDEX idx_consents_type ON user_consents(consent_type);
CREATE INDEX idx_consents_granted ON user_consents(granted);
CREATE INDEX idx_consents_created_at ON user_consents(created_at);

-- ============================================================================
-- LOGIN ACTIVITY (Audit)
-- ============================================================================
CREATE TABLE login_activity (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    success BOOLEAN NOT NULL,
    
    -- Location & device
    ip_address VARCHAR(50) NOT NULL,
    user_agent TEXT,
    location VARCHAR(255),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    device_info JSONB,
    
    -- Details
    reason TEXT, -- For failed attempts
    session_id INT REFERENCES user_sessions(id) ON DELETE SET NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create partitions for login_activity (last 6 months + current month)
CREATE TABLE login_activity_2026_01 PARTITION OF login_activity
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE login_activity_2026_02 PARTITION OF login_activity
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE login_activity_2026_03 PARTITION OF login_activity
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE INDEX idx_login_activity_user_id ON login_activity(user_id);
CREATE INDEX idx_login_activity_email ON login_activity(email);
CREATE INDEX idx_login_activity_created_at ON login_activity(created_at);
CREATE INDEX idx_login_activity_success ON login_activity(success);
CREATE INDEX idx_login_activity_ip_address ON login_activity(ip_address);

-- ============================================================================
-- USER ACTIVITY LOG (General Audit)
-- ============================================================================
CREATE TABLE user_activity_log (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL, -- 'created', 'updated', 'deleted', 'viewed'
    entity VARCHAR(50) NOT NULL, -- 'user', 'appointment', 'prescription'
    entity_id INT,
    
    -- Change tracking
    changes JSONB, -- Old and new values
    
    -- Context
    ip_address VARCHAR(50) NOT NULL,
    user_agent TEXT,
    session_id INT REFERENCES user_sessions(id) ON DELETE SET NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create partitions for user_activity_log
CREATE TABLE user_activity_log_2026_01 PARTITION OF user_activity_log
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE user_activity_log_2026_02 PARTITION OF user_activity_log
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE user_activity_log_2026_03 PARTITION OF user_activity_log
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE INDEX idx_activity_log_user_id ON user_activity_log(user_id);
CREATE INDEX idx_activity_log_action ON user_activity_log(action);
CREATE INDEX idx_activity_log_entity ON user_activity_log(entity);
CREATE INDEX idx_activity_log_created_at ON user_activity_log(created_at);

-- ============================================================================
-- SECURITY EVENTS
-- ============================================================================
CREATE TABLE security_events (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL, -- 'suspicious_login', 'password_reset', 'account_lockout'
    severity severity_type NOT NULL,
    description TEXT NOT NULL,
    
    -- Context
    ip_address VARCHAR(50) NOT NULL,
    user_agent TEXT,
    location VARCHAR(255),
    metadata JSONB,
    
    -- Resolution
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP,
    resolved_by INT REFERENCES users(id) ON DELETE SET NULL,
    resolution_notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_security_events_user_id ON security_events(user_id);
CREATE INDEX idx_security_events_event_type ON security_events(event_type);
CREATE INDEX idx_security_events_severity ON security_events(severity);
CREATE INDEX idx_security_events_resolved ON security_events(resolved);
CREATE INDEX idx_security_events_created_at ON security_events(created_at);

-- ============================================================================
-- ACCOUNT STATUS CHANGES (Audit)
-- ============================================================================
CREATE TABLE account_status_changes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    old_status account_status_type NOT NULL,
    new_status account_status_type NOT NULL,
    reason TEXT,
    changed_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_status_changes_user_id ON account_status_changes(user_id);
CREATE INDEX idx_status_changes_changed_by ON account_status_changes(changed_by);
CREATE INDEX idx_status_changes_created_at ON account_status_changes(created_at);

-- ============================================================================
-- RATE LIMITING
-- ============================================================================
CREATE TABLE rate_limit_rules (
    id SERIAL PRIMARY KEY,
    action VARCHAR(100) UNIQUE NOT NULL,
    max_attempts INT NOT NULL,
    window_size_seconds INT NOT NULL,
    block_duration_seconds INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE rate_limit_log (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL, -- IP, user_id, email
    action VARCHAR(100) NOT NULL,
    count INT DEFAULT 1,
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    blocked BOOLEAN DEFAULT FALSE,
    blocked_until TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rate_limit_identifier ON rate_limit_log(identifier);
CREATE INDEX idx_rate_limit_action ON rate_limit_log(action);
CREATE INDEX idx_rate_limit_window_end ON rate_limit_log(window_end);
CREATE INDEX idx_rate_limit_blocked ON rate_limit_log(blocked);

-- ============================================================================
-- EMAIL QUEUE
-- ============================================================================
CREATE TABLE email_queue (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    to_email VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    
    -- Template support
    template_name VARCHAR(100),
    template_data JSONB,
    
    -- Status
    status email_status DEFAULT 'pending',
    priority INT DEFAULT 5, -- 1 (highest) to 10 (lowest)
    
    -- Retry logic
    attempts INT DEFAULT 0,
    max_attempts INT DEFAULT 3,
    last_attempt_at TIMESTAMP,
    next_attempt_at TIMESTAMP,
    
    -- Result
    error_message TEXT,
    sent_at TIMESTAMP,
    
    -- Tracking
    opened_at TIMESTAMP,
    clicked_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_queue_status ON email_queue(status);
CREATE INDEX idx_email_queue_created_at ON email_queue(created_at);
CREATE INDEX idx_email_queue_next_attempt_at ON email_queue(next_attempt_at);
CREATE INDEX idx_email_queue_priority ON email_queue(priority);

-- ============================================================================
-- API KEYS (For API Access)
-- ============================================================================
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    key_prefix VARCHAR(20) NOT NULL, -- First few chars for identification
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Permissions
    scopes TEXT[], -- Array of permissions: ['read:users', 'write:appointments']
    
    -- Usage
    last_used_at TIMESTAMP,
    usage_count INT DEFAULT 0,
    
    -- Lifecycle
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP,
    revoked_by INT REFERENCES users(id) ON DELETE SET NULL,
    revoke_reason TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_key_prefix ON api_keys(key_prefix);
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at);

-- ============================================================================
-- NOTIFICATIONS
-- ============================================================================
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Content
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'info', 'warning', 'error', 'success'
    category VARCHAR(50), -- 'appointment', 'message', 'system'
    
    -- Action
    action_url TEXT,
    action_text VARCHAR(100),
    
    -- Status
    read_at TIMESTAMP,
    archived_at TIMESTAMP,
    
    -- Related entity
    entity_type VARCHAR(50),
    entity_id INT,
    
    -- Metadata
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_read_at ON notifications(read_at);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_type ON notifications(type);

-- ============================================================================
-- SCHEDULED JOBS
-- ============================================================================
CREATE TABLE scheduled_jobs (
    id SERIAL PRIMARY KEY,
    job_name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    
    -- Schedule
    cron_expression VARCHAR(100),
    last_run_at TIMESTAMP,
    next_run_at TIMESTAMP NOT NULL,
    
    -- Status
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'running', 'completed', 'failed'
    
    -- Result
    error_message TEXT,
    execution_time_ms INT,
    
    -- Control
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_scheduled_jobs_next_run_at ON scheduled_jobs(next_run_at);
CREATE INDEX idx_scheduled_jobs_is_active ON scheduled_jobs(is_active);

-- ============================================================================
-- DATA EXPORT REQUESTS (GDPR)
-- ============================================================================
CREATE TABLE data_export_requests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Request details
    request_type VARCHAR(50) NOT NULL, -- 'export', 'deletion'
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    
    -- Export details
    file_url TEXT,
    file_size_bytes BIGINT,
    expires_at TIMESTAMP,
    
    -- Processing
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    
    -- Audit
    ip_address VARCHAR(50),
    user_agent TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_data_export_user_id ON data_export_requests(user_id);
CREATE INDEX idx_data_export_status ON data_export_requests(status);
CREATE INDEX idx_data_export_created_at ON data_export_requests(created_at);

-- ============================================================================
-- USER TAGS (For Segmentation/Organization)
-- ============================================================================
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7), -- Hex color code
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_tags (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    assigned_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, tag_id)
);

CREATE INDEX idx_user_tags_user_id ON user_tags(user_id);
CREATE INDEX idx_user_tags_tag_id ON user_tags(tag_id);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_credentials_updated_at 
    BEFORE UPDATE ON user_credentials 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_security_info_updated_at 
    BEFORE UPDATE ON user_security_info 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_social_auth_updated_at 
    BEFORE UPDATE ON user_social_auth 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_push_tokens_updated_at 
    BEFORE UPDATE ON user_push_tokens 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_preferences_updated_at 
    BEFORE UPDATE ON user_preferences 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rate_limit_log_updated_at 
    BEFORE UPDATE ON rate_limit_log 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rate_limit_rules_updated_at 
    BEFORE UPDATE ON rate_limit_rules 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_scheduled_jobs_updated_at 
    BEFORE UPDATE ON scheduled_jobs 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Log account status changes
CREATE OR REPLACE FUNCTION log_account_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.account_status IS DISTINCT FROM NEW.account_status THEN
        INSERT INTO account_status_changes (
            user_id, 
            old_status, 
            new_status, 
            reason, 
            changed_by,
            ip_address
        ) VALUES (
            NEW.id, 
            OLD.account_status, 
            NEW.account_status, 
            NEW.status_reason, 
            NEW.status_changed_by,
            inet_client_addr()::VARCHAR
        );
        
        -- Create security event for critical status changes
        IF NEW.account_status IN ('suspended', 'banned', 'closed') THEN
            INSERT INTO security_events (
                user_id,
                event_type,
                severity,
                description,
                ip_address
            ) VALUES (
                NEW.id,
                'account_status_changed',
                'high',
                'Account status changed from ' || OLD.account_status || ' to ' || NEW.account_status,
                inet_client_addr()::VARCHAR
            );
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER track_account_status_changes
    AFTER UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION log_account_status_change();

-- Initialize user preferences on user creation
CREATE OR REPLACE FUNCTION initialize_user_defaults()
RETURNS TRIGGER AS $$
BEGIN
    -- Create security info
    INSERT INTO user_security_info (user_id) 
    VALUES (NEW.id);
    
    -- Create preferences
    INSERT INTO user_preferences (user_id) 
    VALUES (NEW.id);
    
    -- Create credentials record
    INSERT INTO user_credentials (user_id) 
    VALUES (NEW.id);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER initialize_new_user
    AFTER INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION initialize_user_defaults();

-- Update session last_used_at
CREATE OR REPLACE FUNCTION update_session_last_used()
RETURNS TRIGGER AS $$
BEGIN
    -- This would be called by application, but can be automated
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- UTILITY FUNCTIONS
-- ============================================================================

-- Clean expired tokens
CREATE OR REPLACE FUNCTION clean_expired_tokens()
RETURNS void AS $$
BEGIN
    DELETE FROM email_verification_tokens 
    WHERE expires_at < CURRENT_TIMESTAMP AND used_at IS NULL;
    
    DELETE FROM password_reset_tokens 
    WHERE expires_at < CURRENT_TIMESTAMP AND used_at IS NULL;
    
    DELETE FROM phone_verification_codes
    WHERE expires_at < CURRENT_TIMESTAMP AND used_at IS NULL;
    
    DELETE FROM user_sessions 
    WHERE expires_at < CURRENT_TIMESTAMP AND revoked_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- Archive old login activity
CREATE OR REPLACE FUNCTION archive_old_login_activity(months_to_keep INT DEFAULT 6)
RETURNS INT AS $$
DECLARE
    deleted_count INT;
BEGIN
    DELETE FROM login_activity 
    WHERE created_at < CURRENT_TIMESTAMP - (months_to_keep || ' months')::INTERVAL;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Archive old activity logs
CREATE OR REPLACE FUNCTION archive_old_activity_logs(months_to_keep INT DEFAULT 6)
RETURNS INT AS $$
DECLARE
    deleted_count INT;
BEGIN
    DELETE FROM user_activity_log 
    WHERE created_at < CURRENT_TIMESTAMP - (months_to_keep || ' months')::INTERVAL;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Soft delete user
CREATE OR REPLACE FUNCTION soft_delete_user(
    user_id_param INT, 
    deleted_by_param INT,
    reason_param TEXT DEFAULT NULL
)
RETURNS void AS $$
BEGIN
    -- Update user
    UPDATE users 
    SET 
        deleted_at = CURRENT_TIMESTAMP,
        deleted_by = deleted_by_param,
        account_status = 'closed',
        status_reason = reason_param,
        status_changed_at = CURRENT_TIMESTAMP,
        status_changed_by = deleted_by_param
    WHERE id = user_id_param AND deleted_at IS NULL;
    
    -- Revoke all sessions
    UPDATE user_sessions 
    SET 
        revoked_at = CURRENT_TIMESTAMP,
        revoked_by = deleted_by_param,
        revoke_reason = 'User account deleted'
    WHERE user_id = user_id_param AND revoked_at IS NULL;
    
    -- Deactivate push tokens
    UPDATE user_push_tokens 
    SET 
        is_active = FALSE,
        deactivated_at = CURRENT_TIMESTAMP,
        deactivation_reason = 'User account deleted'
    WHERE user_id = user_id_param AND is_active = TRUE;
    
    -- Revoke API keys
    UPDATE api_keys
    SET 
        revoked_at = CURRENT_TIMESTAMP,
        revoked_by = deleted_by_param,
        revoke_reason = 'User account deleted'
    WHERE user_id = user_id_param AND revoked_at IS NULL;
    
    -- Create security event
    INSERT INTO security_events (
        user_id,
        event_type,
        severity,
        description,
        ip_address
    ) VALUES (
        user_id_param,
        'account_deleted',
        'high',
        'User account was soft deleted. Reason: ' || COALESCE(reason_param, 'Not specified'),
        inet_client_addr()::VARCHAR
    );
END;
$$ LANGUAGE plpgsql;

-- Restore user
CREATE OR REPLACE FUNCTION restore_user(user_id_param INT)
RETURNS void AS $$
BEGIN
    UPDATE users 
    SET 
        deleted_at = NULL,
        deleted_by = NULL,
        account_status = 'active'
    WHERE id = user_id_param AND deleted_at IS NOT NULL;
    
    -- Create security event
    INSERT INTO security_events (
        user_id,
        event_type,
        severity,
        description,
        ip_address
    ) VALUES (
        user_id_param,
        'account_restored',
        'medium',
        'User account was restored',
        inet_client_addr()::VARCHAR
    );
END;
$$ LANGUAGE plpgsql;

-- Check rate limit
CREATE OR REPLACE FUNCTION check_rate_limit(
    identifier_param VARCHAR(255),
    action_param VARCHAR(100)
) RETURNS BOOLEAN AS $$
DECLARE
    rule_record RECORD;
    current_count INT;
    existing_record RECORD;
BEGIN
    -- Get the rule
    SELECT * INTO rule_record 
    FROM rate_limit_rules 
    WHERE action = action_param AND is_active = TRUE;
    
    IF NOT FOUND THEN
        RETURN TRUE; -- No rule, allow
    END IF;
    
    -- Check for existing window
    SELECT * INTO existing_record 
    FROM rate_limit_log
    WHERE identifier = identifier_param 
    AND action = action_param
    AND window_end > CURRENT_TIMESTAMP
    ORDER BY window_end DESC
    LIMIT 1;
    
    -- Check if blocked
    IF existing_record.blocked AND existing_record.blocked_until > CURRENT_TIMESTAMP THEN
        RETURN FALSE;
    END IF;
    
    IF FOUND AND existing_record.window_end > CURRENT_TIMESTAMP THEN
        -- Update existing window
        current_count := existing_record.count + 1;
        
        UPDATE rate_limit_log
        SET 
            count = current_count,
            blocked = (current_count >= rule_record.max_attempts),
            blocked_until = CASE 
                WHEN current_count >= rule_record.max_attempts 
                THEN CURRENT_TIMESTAMP + (rule_record.block_duration_seconds || ' seconds')::INTERVAL
                ELSE NULL
            END,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = existing_record.id;
        
        RETURN current_count < rule_record.max_attempts;
    ELSE
        -- Create new window
        INSERT INTO rate_limit_log (
            identifier, 
            action, 
            count, 
            window_start, 
            window_end
        ) VALUES (
            identifier_param, 
            action_param, 
            1,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP + (rule_record.window_size_seconds || ' seconds')::INTERVAL
        );
        
        RETURN TRUE;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Search users (Full-text search)
CREATE OR REPLACE FUNCTION search_users(
    search_query TEXT,
    limit_param INT DEFAULT 50
)
RETURNS TABLE (
    id INT,
    email VARCHAR,
    name VARCHAR,
    role user_role,
    rank REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.email,
        u.name,
        u.role,
        ts_rank(
            to_tsvector('english', u.name || ' ' || u.email),
            plainto_tsquery('english', search_query)
        ) AS rank
    FROM users u
    WHERE u.deleted_at IS NULL
    AND u.account_status = 'active'
    AND (
        to_tsvector('english', u.name || ' ' || u.email) @@ plainto_tsquery('english', search_query)
    )
    ORDER BY rank DESC
    LIMIT limit_param;
END;
$$ LANGUAGE plpgsql;

-- Get user security summary
CREATE OR REPLACE FUNCTION get_user_security_summary(user_id_param INT)
RETURNS TABLE (
    two_factor_enabled BOOLEAN,
    password_age_days INT,
    active_sessions_count INT,
    recent_failed_logins INT,
    is_locked BOOLEAN,
    security_score INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        uc.two_factor_enabled,
        EXTRACT(DAY FROM CURRENT_TIMESTAMP - usi.last_password_change)::INT AS password_age_days,
        (SELECT COUNT(*) FROM user_sessions WHERE user_id = user_id_param AND revoked_at IS NULL AND expires_at > CURRENT_TIMESTAMP)::INT AS active_sessions_count,
        (SELECT COUNT(*) FROM login_activity WHERE user_id = user_id_param AND success = FALSE AND created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours')::INT AS recent_failed_logins,
        (usi.locked_until IS NOT NULL AND usi.locked_until > CURRENT_TIMESTAMP) AS is_locked,
        usi.security_score
    FROM user_credentials uc
    JOIN user_security_info usi ON usi.user_id = uc.user_id
    WHERE uc.user_id = user_id_param;
END;
$$ LANGUAGE plpgsql;

-- Calculate and update user security score
CREATE OR REPLACE FUNCTION calculate_security_score(user_id_param INT)
RETURNS INT AS $$
DECLARE
    score INT := 100;
    has_2fa BOOLEAN;
    password_age_days INT;
    recent_failed_logins INT;
    active_sessions INT;
BEGIN
    -- Check 2FA (+20 if enabled)
    SELECT two_factor_enabled INTO has_2fa 
    FROM user_credentials WHERE user_id = user_id_param;
    
    IF NOT has_2fa THEN
        score := score - 20;
    END IF;
    
    -- Check password age (-10 if > 90 days)
    SELECT EXTRACT(DAY FROM CURRENT_TIMESTAMP - last_password_change)::INT 
    INTO password_age_days
    FROM user_security_info WHERE user_id = user_id_param;
    
    IF password_age_days > 90 THEN
        score := score - 10;
    END IF;
    
    -- Check recent failed logins (-5 per failed login in last 24h)
    SELECT COUNT(*) INTO recent_failed_logins
    FROM login_activity 
    WHERE user_id = user_id_param 
    AND success = FALSE 
    AND created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours';
    
    score := score - (recent_failed_logins * 5);
    
    -- Check active sessions (-5 if > 5 active sessions)
    SELECT COUNT(*) INTO active_sessions
    FROM user_sessions 
    WHERE user_id = user_id_param 
    AND revoked_at IS NULL 
    AND expires_at > CURRENT_TIMESTAMP;
    
    IF active_sessions > 5 THEN
        score := score - 10;
    END IF;
    
    -- Ensure score is between 0 and 100
    score := GREATEST(0, LEAST(100, score));
    
    -- Update the score
    UPDATE user_security_info 
    SET 
        security_score = score,
        security_score_updated_at = CURRENT_TIMESTAMP
    WHERE user_id = user_id_param;
    
    RETURN score;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- INITIAL DATA
-- ============================================================================

-- Insert default rate limit rules
INSERT INTO rate_limit_rules (action, max_attempts, window_size_seconds, block_duration_seconds, description) VALUES
('login', 5, 300, 1800, '5 login attempts per 5 minutes, block for 30 minutes'),
('password_reset', 3, 3600, 3600, '3 password reset attempts per hour, block for 1 hour'),
('email_verify', 5, 3600, 3600, '5 email verification attempts per hour, block for 1 hour'),
('resend_email', 3, 600, 1800, '3 email resend attempts per 10 minutes, block for 30 minutes'),
('api_request', 1000, 3600, 3600, '1000 API requests per hour, block for 1 hour'),
('password_change', 3, 86400, 86400, '3 password changes per day, block for 24 hours');

-- Insert default scheduled jobs
INSERT INTO scheduled_jobs (job_name, description, next_run_at) VALUES
('clean_expired_tokens', 'Clean expired tokens from database', CURRENT_TIMESTAMP),
('archive_old_logs', 'Archive old login activity and audit logs', CURRENT_TIMESTAMP + INTERVAL '1 day'),
('process_email_queue', 'Process pending emails in queue', CURRENT_TIMESTAMP),
('calculate_security_scores', 'Recalculate security scores for all users', CURRENT_TIMESTAMP + INTERVAL '1 hour'),
('cleanup_revoked_sessions', 'Remove old revoked sessions', CURRENT_TIMESTAMP + INTERVAL '7 days');

-- Insert default tags
INSERT INTO tags (name, description, color) VALUES
('vip', 'VIP users', '#FFD700'),
('premium', 'Premium subscription users', '#9B59B6'),
('verified', 'Verified users', '#3498DB'),
('staff', 'Staff members', '#E74C3C'),
('beta_tester', 'Beta testers', '#F39C12');

-- ============================================================================
-- VIEWS (Helpful for queries)
-- ============================================================================

-- Active users view
CREATE OR REPLACE VIEW active_users AS
SELECT 
    u.id,
    u.email,
    u.name,
    u.role,
    u.account_status,
    u.email_verified,
    u.phone_verified,
    u.last_active_at,
    u.created_at,
    uc.two_factor_enabled,
    usi.last_login_at,
    usi.security_score
FROM users u
LEFT JOIN user_credentials uc ON u.id = uc.user_id
LEFT JOIN user_security_info usi ON u.id = usi.user_id
WHERE u.deleted_at IS NULL
AND u.account_status = 'active';

-- User sessions view
CREATE OR REPLACE VIEW active_sessions AS
SELECT 
    s.id,
    s.user_id,
    u.email,
    u.name,
    s.device_name,
    s.platform,
    s.ip_address,
    s.location,
    s.last_used_at,
    s.created_at,
    s.expires_at
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.revoked_at IS NULL 
AND s.expires_at > CURRENT_TIMESTAMP;

-- Security events summary view
CREATE OR REPLACE VIEW recent_security_events AS
SELECT 
    se.id,
    se.user_id,
    u.email,
    u.name,
    se.event_type,
    se.severity,
    se.description,
    se.resolved,
    se.created_at
FROM security_events se
JOIN users u ON se.user_id = u.id
WHERE se.created_at > CURRENT_TIMESTAMP - INTERVAL '7 days'
ORDER BY se.created_at DESC;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE users IS 'Core user entity with soft delete support and comprehensive profile management';
COMMENT ON TABLE user_credentials IS 'Separated credentials for enhanced security';
COMMENT ON TABLE user_security_info IS 'Security metadata separated for performance optimization';
COMMENT ON TABLE user_sessions IS 'User session management with device tracking';
COMMENT ON TABLE user_push_tokens IS 'Push notification tokens - one user can have multiple devices';
COMMENT ON TABLE user_preferences IS 'User preferences for localization, notifications, and display';
COMMENT ON TABLE user_consents IS 'GDPR compliance - track user consents and data processing agreements';
COMMENT ON TABLE rate_limit_log IS 'Rate limiting for abuse prevention';
COMMENT ON TABLE email_queue IS 'Reliable email delivery with retry logic';
COMMENT ON TABLE api_keys IS 'API key management for programmatic access';
COMMENT ON TABLE notifications IS 'In-app notifications system';
COMMENT ON TABLE data_export_requests IS 'GDPR data export and deletion requests';

COMMENT ON COLUMN password_reset_tokens.used_by_session_id IS 'Track which session used the token';
COMMENT ON COLUMN email_verification_tokens.used_by_session_id IS 'Track which session used the token';
COMMENT ON COLUMN user_security_info.security_score IS 'Security score from 0-100 based on multiple factors';

-- ============================================================================
-- GRANTS (Adjust based on your user roles)
-- ============================================================================

-- Example: Grant permissions to application user
-- GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO app_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
-- GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO app_user;

-- ============================================================================
-- END OF SCHEMA
-- ============================================================================

-- Print success message
DO $$
BEGIN
    RAISE NOTICE 'User management schema created successfully!';
    RAISE NOTICE 'Tables created: 30+';
    RAISE NOTICE 'Indexes created: 100+';
    RAISE NOTICE 'Functions created: 10+';
    RAISE NOTICE 'Triggers created: 10+';
    RAISE NOTICE 'Views created: 3';
END $$;