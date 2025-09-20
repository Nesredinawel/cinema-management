🚀 Cinema Auth & Scheduling - Docker Setup

This project uses Docker Compose to spin up a full environment with:

Postgres – database for authentication and scheduling

Redis – caching and session management

Hasura – GraphQL engine on top of Postgres

Auth Backend – custom authentication service

Cinema Scheduling – scheduling service

📦 Services Overview
1. Postgres

Runs a PostgreSQL database.

Stores authentication and scheduling data.

Persists data with a named volume.

2. Redis

Provides caching and session storage.

Accessible on port 6379.

3. Hasura

GraphQL API layer over Postgres.

Admin console enabled.

Depends on Postgres.

4. Auth Backend

Custom backend for authentication.

Connects to Postgres and Redis.

Exposes APIs on port 8081.

5. Cinema Scheduling

Scheduling microservice.

Uses the same database as the auth service.

Exposes APIs on port 8082.

⚙️ Environment Variables

For security reasons, do not hardcode credentials directly in docker-compose.yaml.
Instead, create a .env file in your project root with values like:

# Database
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=cinema_auth
DB_HOST=postgres
DB_PORT=5432

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Hasura
HASURA_GRAPHQL_DATABASE_URL=postgres://your_db_user:your_db_password@postgres:5432/cinema_auth
HASURA_GRAPHQL_ADMIN_SECRET=your_admin_secret
HASURA_GRAPHQL_ENABLE_CONSOLE=true

# Auth Backend
PORT=8081
JWT_SECRET=your_jwt_secret

# Cinema Scheduling
PORT=8082
JWT_SECRET=your_jwt_secret


👉 Your docker-compose.yaml should then reference these variables like ${DB_USER}, ${DB_PASSWORD}, etc.

▶️ Getting Started



Create .env file for each sub folders
Copy the example and update values:
# ==============================
# 🌐 General
# ==============================
PORT=8082
JWT_SECRET=your_jwt_secret_here

# ==============================
# 🛢️ Postgres Database
# ==============================
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=your_password_here
POSTGRES_DB=cinema_auth

# ==============================
# 🔑 OAuth (Google Example)
# ==============================
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8082/auth/google/callback

cp .env.example .env


Start all services

docker-compose up -d --build


Check running containers

docker ps

🌍 Service Endpoints

Postgres → localhost:5432

Redis → localhost:6379

Hasura Console → http://localhost:8080

Auth Backend API → http://localhost:8081

Cinema Scheduling API → http://localhost:8082

🗑️ Stopping & Cleaning Up

To stop services:

docker-compose down


To remove containers, networks, and volumes:

docker-compose down -v
