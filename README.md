# 🚀 Cinema Auth & Scheduling - Docker Setup

This project uses **Docker Compose** to spin up a full environment with:

* **Postgres** – database for authentication and scheduling
* **Redis** – caching and session management
* **Hasura** – GraphQL engine on top of Postgres
* **Auth Backend** – custom authentication service
* **Cinema Scheduling** – scheduling service

---

## 📦 Services Overview

### 1. **Postgres**

* Runs a PostgreSQL database.
* Stores authentication and scheduling data.
* Persists data with a named volume.

### 2. **Redis**

* Provides caching and session storage.
* Accessible on port **6379**.

### 3. **Hasura**

* GraphQL API layer over Postgres.
* Admin console enabled.
* Depends on **Postgres**.

### 4. **Auth Backend**

* Custom backend for authentication.
* Connects to **Postgres** and **Redis**.
* Exposes APIs on port **8081**.

### 5. **Cinema Scheduling**

* Scheduling microservice.
* Uses the same database as the auth service.
* Exposes APIs on port **8082**.

---

## ⚙️ Environment Variables

For security reasons, **do not hardcode credentials** directly in `docker-compose.yaml`.
Instead, create a `.env` file in your project root and subfolders with values like:

```dotenv
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
```

Copy the example and create your own `.env`:

```bash
cp .env.example .env
```

---

## 🐳 Example `docker-compose.yaml`

Here’s a reference `docker-compose.yaml` that uses environment variables:

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:15
    container_name: auth_postgres
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - auth_postgres_data:/var/lib/postgresql/data
    restart: always

  redis:
    image: redis:7
    container_name: auth_redis
    ports:
      - "6379:6379"
    restart: always

  hasura:
    image: hasura/graphql-engine:v2.3.0
    container_name: auth_hasura
    depends_on:
      - postgres
    ports:
      - "8080:8080"
    environment:
      HASURA_GRAPHQL_DATABASE_URL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:${POSTGRES_PORT}/${POSTGRES_DB}
      HASURA_GRAPHQL_ADMIN_SECRET: ${HASURA_GRAPHQL_ADMIN_SECRET}
      HASURA_GRAPHQL_ENABLE_CONSOLE: "true"
    restart: always

  auth-backend:
    build: ./auth-backend
    container_name: cinema_auth_backend
    depends_on:
      - postgres
      - redis
      - hasura
    ports:
      - "8081:8081"
    environment:
      PORT: 8081
      DB_HOST: ${POSTGRES_HOST}
      DB_PORT: ${POSTGRES_PORT}
      DB_USER: ${POSTGRES_USER}
      DB_PASSWORD: ${POSTGRES_PASSWORD}
      DB_NAME: ${POSTGRES_DB}
      REDIS_HOST: ${REDIS_HOST}
      REDIS_PORT: ${REDIS_PORT}
      JWT_SECRET: ${JWT_SECRET}
    restart: always

  cinema-scheduling:
    build: ./cinema-scheduling
    container_name: cinema_scheduling_service
    depends_on:
      - auth-backend
    ports:
      - "8082:8082"
    environment:
      PORT: 8082
      DB_HOST: ${POSTGRES_HOST}
      DB_PORT: ${POSTGRES_PORT}
      DB_USER: ${POSTGRES_USER}
      DB_PASSWORD: ${POSTGRES_PASSWORD}
      DB_NAME: ${POSTGRES_DB}
      JWT_SECRET: ${JWT_SECRET}
    restart: always

volumes:
  auth_postgres_data:
```

---

## ▶️ Getting Started

1. **Create `.env` file**

   ```bash
   cp .env.example .env
   ```

2. **Start all services**

   ```bash
   docker-compose up -d --build
   ```

3. **Check running containers**

   ```bash
   docker ps
   ```

---

## 🌍 Service Endpoints

* **Postgres** → `localhost:5432`
* **Redis** → `localhost:6379`
* **Hasura Console** → [http://localhost:8080](http://localhost:8080)
* **Auth Backend API** → [http://localhost:8081](http://localhost:8081)
* **Cinema Scheduling API** → [http://localhost:8082](http://localhost:8082)

---

## 🗑️ Stopping & Cleaning Up

Stop services:

```bash
docker-compose down
```

Remove containers, networks, and volumes:

```bash
docker-compose down -v
```

---


