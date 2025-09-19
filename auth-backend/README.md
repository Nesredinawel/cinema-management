# Auth Backend

This folder contains the authentication backend service for the Cinema Management system.

## Features

- User registration and login
- JWT-based authentication
- Password hashing and validation
- Role-based access control

## Getting Started

1. **Install dependencies:**
    ```bash
    npm install
    ```
2. **Configure environment variables:**  
    Copy `.env.example` to `.env` and update the values as needed.

3. **Run the server:**
    ```bash
    npm start
    ```

## Folder Structure

- `src/` - Source code for the authentication service
- `tests/` - Unit and integration tests
- `config/` - Configuration files

## API Endpoints

- `POST /register` - Register a new user
- `POST /login` - Authenticate a user and receive a token
- `GET /profile` - Get user profile (requires authentication)

## License

This project is licensed under the MIT License.
