 
-- DROP TABLE IF EXISTS notifications;
-- DROP TABLE IF EXISTS survival_guides;
-- DROP TABLE IF EXISTS user_location;
-- DROP TABLE IF EXISTS checklist_items;
-- DROP TABLE IF EXISTS users;

-- USER
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT, 
    avatar TEXT, 
    push_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CHECKLIST
CREATE TABLE IF NOT EXISTS checklist_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,           
    category VARCHAR(100),         
    Description TEXT NULL,
    is_checked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_checklist_user ON checklist_items(user_id, created_at DESC);
 
-- USER LOCATION
CREATE TABLE IF NOT EXISTS user_location (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    lat DOUBLE PRECISION, 
    lon DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SURVIVAL GUIDE
CREATE TABLE IF NOT EXISTS survival_guides (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title TEXT NOT NULL,
    category TEXT,
    difficulty TEXT,
    icon TEXT,
    content TEXT NOT NULL,
    image_url TEXT,
    views INTEGER DEFAULT 0,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- NOTIFICATIONS
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY, 
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB, -- Deep link, action data
    type VARCHAR(50) DEFAULT 'in_app', -- 'in_app', 'push', 'both'
    is_read BOOLEAN DEFAULT FALSE,
    sent BOOLEAN DEFAULT FALSE, -- Cho push notification
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created 
ON notifications(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_notifications_unread 
ON notifications(user_id, is_read) WHERE is_read = FALSE;

CREATE INDEX IF NOT EXISTS idx_notifications_unsent_push 
ON notifications(user_id, sent) WHERE sent = FALSE AND type IN ('push', 'both');