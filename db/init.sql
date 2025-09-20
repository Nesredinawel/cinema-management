CREATE DATABASE cinema_auth;
CREATE DATABASE cinema_scheduling;


-- Switch to cinema_auth database
\c cinema_auth;

-- ---------------- Roles Table ----------------
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL -- "admin", "staff", "customer"
);

-- ---------------- Users Table ----------------
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT UNIQUE,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    google_id TEXT UNIQUE,
    is_verified BOOLEAN DEFAULT FALSE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ---------------- OTP History Table ----------------
CREATE TABLE IF NOT EXISTS otp_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    phone TEXT NOT NULL,
    code TEXT NOT NULL,
    status TEXT NOT NULL, -- "SENT", "VERIFIED", "FAILED", "EXPIRED"
    failed_attempts INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    verified_at TIMESTAMP
);

-- ---------------- Admin Roles Table ----------------
CREATE TABLE IF NOT EXISTS admin_roles (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    level TEXT NOT NULL -- "admin", "manager"
);

-- ---------------- Staff Roles Table ----------------
CREATE TABLE IF NOT EXISTS staff_roles (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    dept TEXT NOT NULL
);

-- ---------------- Customer Roles Table ----------------
CREATE TABLE IF NOT EXISTS customer_roles (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    loyalty_points INT DEFAULT 0
);

-- ---------------- Seed roles ----------------
INSERT INTO roles (name) VALUES 
    ('admin') 
    ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (name) VALUES 
    ('staff') 
    ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (name) VALUES 
    ('customer') 
    ON CONFLICT (name) DO NOTHING;






-- Switch to cinema_scheduling database
\c cinema_scheduling;

-- ==============================
-- Table: genres
-- ==============================
CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ==============================
-- Table: halls
-- ==============================
CREATE TABLE IF NOT EXISTS halls (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    location VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ==============================
-- Table: movies
-- ==============================
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    duration INT NOT NULL,
    release_year INT NOT NULL,
    rating NUMERIC(3,1),
    image_poster_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ==============================
-- Table: movie_genres (many-to-many)
-- ==============================
CREATE TABLE IF NOT EXISTS movie_genres (
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    genre_id INT NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (movie_id, genre_id)
);

-- ==============================
-- Table: schedules
-- ==============================
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    hall_id INT NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
    show_time TIMESTAMP NOT NULL,
    available_seats INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ==============================
-- Table: snacks
-- ==============================
CREATE TABLE IF NOT EXISTS snacks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    description TEXT,
    category VARCHAR(255),
    snack_image_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ==============================
-- Table: schedule_snacks (snacks assigned to a schedule)
-- ==============================
CREATE TABLE IF NOT EXISTS schedule_snacks (
    id SERIAL PRIMARY KEY,
    schedule_id INT NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    snack_id INT NOT NULL REFERENCES snacks(id) ON DELETE CASCADE,
    available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (schedule_id, snack_id)
);
