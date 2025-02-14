# Users Service

## Overview
The **Users Service** is a Golang-based API that manages user accounts, profile images, and club associations. It is built with `Gin Gonic` for the HTTP router and `sqlc` for efficient database interaction with PostgreSQL. The service supports CRUD operations for users and clubs, along with soft deletion and transaction management.

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
|--------|--------------------------------|-----------------------------|
| `POST` | `/api/v1/users/add` | Create a new user |
| `PUT`  | `/api/v1/users/edit/:user_id` | Update user details |
| `PUT`  | `/api/v1/users/edit-img/:user_id` | Update user profile image |
| `DELETE` | `/api/v1/users/delete/:user_id` | Soft delete a user |
| `GET`  | `/api/v1/users/get/:user_id` | Get user by ID |
| `GET`  | `/api/v1/users/get-avatar/:nickname` | Get user avatar by nickname |
| `GET`  | `/api/v1/users/find` | Retrieve all users |
| `GET`  | `/api/v1/users/list` | List users with pagination |

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
curl -X POST http://localhost:8080/api/v1/users/add -H "Content-Type: application/json" -d '{
    "nickname": "john_doe",
    "country": "USA",
    "city": "New York"
}'
```

### Get User by ID
```sh
curl -X GET http://localhost:8080/api/v1/users/get/{user_id}
```

### Get User Avatar by Nickname
```sh
curl -X GET http://localhost:8080/api/v1/users/get-avatar/{nickname}
```

### Find Users
```sh
curl -X GET http://localhost:8080/api/v1/users/find
```

### List Users
```sh
curl -X GET http://localhost:8080/api/v1/users/list
```

### Update User
```sh
curl -X PUT http://localhost:8080/api/v1/users/edit/{user_id} -H "Content-Type: application/json" -d '{
    "city": "Los Angeles"
}'
```

### Update Profile Image
```sh
curl -X PUT http://localhost:8080/api/v1/users/edit-img/{user_id} -H "Content-Type: application/json" -d '{
    "img": "https://example.com/profile.jpg"
}'
```

### Delete User (Soft Delete)
```sh
curl -X DELETE http://localhost:8080/api/v1/users/delete/{user_id}
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
