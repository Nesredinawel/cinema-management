-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255), -- nullable for Google OAuth users
    role VARCHAR(50) NOT NULL DEFAULT 'customer', -- 'customer', 'staff', 'admin'
    profile_image_url VARCHAR(512),
    google_id VARCHAR(255), -- stores Google OAuth ID for customers
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Trigger to automatically update 'updated_at' column on row update
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Optional: Insert an admin user (replace password_hash with bcrypt hash)
INSERT INTO users (name, email, password_hash, role)
VALUES ('Admin User', 'admin@cinema.com', '$2a$12$yourbcryptpasswordhash', 'admin')
ON CONFLICT (email) DO NOTHING;
