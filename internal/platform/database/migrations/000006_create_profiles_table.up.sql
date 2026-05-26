CREATE TABLE IF NOT EXISTS
    profiles (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        full_name VARCHAR(100) NOT NULL,
        avatar_url VARCHAR(2048),
        headline VARCHAR(255),
        bio TEXT,
        resume_url VARCHAR(2048),
        social_links JSONB DEFAULT '{}',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_profiles_created_at ON profiles (created_at DESC);