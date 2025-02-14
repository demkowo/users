# Users Service

## Overview
The **Users Service** is responsible for managing user accounts, including creation, retrieval, updates, and deletion. It also handles user profile images and associations with clubs. The service is implemented in Golang using `sqlc` for efficient database interaction with PostgreSQL.

## Features
- **User Management**: Add, find, list, update, and delete users.
- **Profile Image Handling**: Retrieve and update user profile images.
- **Club Association**: Manage user-club relationships.
- **Soft Deletion**: Users are soft-deleted to preserve data integrity.
- **Transaction Management**: Ensures atomicity in operations.

## Directory Structure
```
users/
│-- models/       # Contains data models
│-- repositories/ # Data access layer (PostgreSQL implementation)
│   ├── postgres/
│   │   ├── sqlc/ # Auto-generated SQL queries
│   │   ├── repository.go # User repository implementation
│-- main.go       # Service entry point
```

## API Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| `POST` | `/users` | Create a new user |
| `GET`  | `/users` | Retrieve all users |
| `GET`  | `/users/{id}` | Get user by ID |
| `PATCH` | `/users/{id}` | Update user details |
| `PATCH` | `/users/{id}/image` | Update user profile image |
| `DELETE` | `/users/{id}` | Soft delete a user |

## Database Schema
The service interacts with the following tables:

### `users`
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    nickname TEXT UNIQUE NOT NULL,
    img TEXT,
    country TEXT,
    city TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);
```

### `clubs`
```sql
CREATE TABLE clubs (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
```

### `user_clubs`
```sql
CREATE TABLE user_clubs (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    club_id UUID REFERENCES clubs(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, club_id)
);
```

## Usage

### Create a User
```sh
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{
    "nickname": "john_doe",
    "country": "USA",
    "city": "New York"
}'
```

### Get All Users
```sh
curl -X GET http://localhost:8080/users
```

### Get User by ID
```sh
curl -X GET http://localhost:8080/users/{id}
```

### Update User
```sh
curl -X PATCH http://localhost:8080/users/{id} -H "Content-Type: application/json" -d '{
    "city": "Los Angeles"
}'
```

### Update Profile Image
```sh
curl -X PATCH http://localhost:8080/users/{id}/image -H "Content-Type: application/json" -d '{
    "img": "https://example.com/profile.jpg"
}'
```

### Delete User (Soft Delete)
```sh
curl -X DELETE http://localhost:8080/users/{id}
```

## Transactions & Error Handling
- All **write operations** (`Add`, `Update`, `Delete`) use transactions to ensure atomicity.
- **Soft deletion** is implemented to prevent accidental data loss.
- Errors are handled gracefully, returning appropriate HTTP status codes.

## Development Setup
### Prerequisites
- Golang (>=1.18)
- PostgreSQL
- `sqlc`

### Install Dependencies
```sh
go mod tidy
```

### Run Service
```sh
go run main.go
```