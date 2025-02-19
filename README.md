# Users Service

## Overview
The **Users Service** is a Golang-based API for managing user accounts, profile images, and club associations. It provides dual API support:
- **REST API** (using Gin Gonic) on port **5000**
- **gRPC API** on port **50000**

The service leverages `sqlc` for efficient PostgreSQL interactions.

## Features
- **User Management**: Create, update, retrieve, and soft-delete users.
- **Profile Image Handling**: Update and fetch user profile images.
- **Club Association**: Manage many-to-many relationships between users and clubs.
- **Dual API Support**: Access via REST or gRPC.
- **Transaction Management**: Ensures atomic operations.

## Directory Structure
```
. 
├── cmd 
│ └── users 
│     └── main.go 
├── internal 
│   ├── app 
│   │   ├── app.go 
│   │   ├── pb_server.go      # gRPC server entry point 
│   │   └── userRoutes.go     # Gin Gonic routes 
│   ├── config 
│   │   └── config.go 
│   ├── generated 
│   │   ├── users_grpc.pb.go  # gRPC generated code 
│   │   └── users.pb.go       # gRPC generated code
│   ├── handlers 
│   │   ├── gin 
│   │   │   └── users_handler.go 
│   │   └── grpc 
│   │       └── users_handler.go 
│   ├── models 
│   │   └── users_model.go 
│   ├── repositories 
│   │   └── postgres 
│   │       ├── queries.sql
│   │       ├── schema.sql 
│   │       ├── sqlc 
│   │       │   ├── db.go 
│   │       │   ├── models.go 
│   │       │   └── queries.sql.go 
│   │       ├── sqlc.yaml 
│   │       └── users_repository.go 
│   └── services 
│   └── users_service.go 
├── proto 
│   └── users.proto 
└── README.md
```

## API Endpoints

### REST API (Port: 5000)
| Method   | Endpoint                                    | Description                            |
|----------|---------------------------------------------|----------------------------------------|
| `POST`   | `/api/v1/users/add`                         | Create a new user                      |
| `PUT`    | `/api/v1/users/edit/:user_id`               | Update user details                    |
| `PUT`    | `/api/v1/users/edit-img/:user_id`           | Update user profile image              |
| `DELETE` | `/api/v1/users/delete/:user_id`             | Soft-delete a user                     |
| `GET`    | `/api/v1/users/get/:user_id`                | Retrieve a user by ID                  |
| `GET`    | `/api/v1/users/get-avatar/:nickname`        | Retrieve user avatar by nickname       |
| `GET`    | `/api/v1/users/find`                        | Retrieve all users                     |
| `GET`    | `/api/v1/users/list`                        | List users with pagination             |

### gRPC API (Port: 50000)
Refer to the [proto/users.proto](proto/users.proto) file for detailed service and message definitions. The gRPC API supports similar operations:
- **AddUser**
- **DeleteUser**
- **FindUsers**
- **GetAvatarByNickname**
- **GetUserById**
- **ListUsers**
- **UpdateUser**
- **UpdateUserImg**

Example using `grpcurl`:
```sh
grpcurl -plaintext -d '{
  "nickname": "john_doe",
  "img": "https://example.com/avatar.jpg",
  "name": "John",
  "last_name": "Doe",
  "country": "USA",
  "city": "New York",
  "clubs": ["Club1", "Club2"]
}' localhost:50000 users.UserService/AddUser
```

## Database Schema
The service interacts with the following tables:

### `users`
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    nickname TEXT UNIQUE NOT NULL,
    img TEXT,
    country TEXT,
    city TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted BOOLEAN DEFAULT false
);
```

### `clubs`
```sql
CREATE TABLE IF NOT EXISTS clubs (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
```

### `user_clubs`
```sql
CREATE TABLE IF NOT EXISTS user_clubs (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    club_id UUID NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, club_id)
);
```

## Usage

### REST API Examples

#### Create a User
```sh
curl -X POST http://localhost:5000/api/v1/users/add \
-H "Content-Type: application/json" \
-d '{
    "nickname": "john_doe",
    "img": "https://example.com/avatar.jpg",
    "name": "John",
    "last_name": "Doe",
    "country": "USA",
    "city": "New York",
    "clubs": ["Club1", "Club2"]
}'
```

#### Get User by ID
```sh
curl -X GET http://localhost:5000/api/v1/users/get/{user_id}
```

#### Get User Avatar by Nickname
```sh
curl -X GET http://localhost:5000/api/v1/users/get-avatar/{nickname}
```

#### Find Users
```sh
curl -X GET http://localhost:5000/api/v1/users/find
```

#### List Users
```sh
curl -X GET "http://localhost:5000/api/v1/users/list?limit=10&offset=0"
```

#### Update User
```sh
curl -X PUT http://localhost:5000/api/v1/users/edit/{user_id} \
-H "Content-Type: application/json" \
-d '{
    "country": "USA",
    "city": "Los Angeles",
    "clubs": ["Club3", "Club4"]
}'
```

#### Update Profile Image
```sh
curl -X PUT http://localhost:5000/api/v1/users/edit-img/{user_id} \
-H "Content-Type: application/json" \
-d '{
    "img": "https://example.com/new_avatar.jpg"
}'
```

#### Delete User (Soft Delete)
```sh
curl -X DELETE http://localhost:5000/api/v1/users/delete/{user_id}
```

### gRPC API

Use a gRPC client (e.g., Postman, grpcurl) to call methods defined in the `users.proto` file on port **50000**.

## Transactions & Error Handling
- All **write operations** (`Add`, `Update`, `Delete`) use transactions to ensure atomicity.
- **Soft deletion** is implemented to prevent accidental data loss.
- Errors are handled gracefully, returning appropriate HTTP status codes.

## Development Setup
### Prerequisites
- Golang (>=1.18)
- PostgreSQL
- `sqlc`
- Protocol Buffers (protoc) & gRPC

### Install Dependencies
```sh
go mod tidy
```

### Generate Code

#### gRPC

use `make_proto_files.sh` or call:

```sh
protoc \
  -I=proto \
  --go_out=internal/generated \
  --go_opt=paths=source_relative \
  --go-grpc_out=internal/generated \
  --go-grpc_opt=paths=source_relative \
  proto/users.proto
```

#### sqlc files
```sh
cd internal/repositories/postgres && sqlc generate
```

### Run Service
```sh
go run main.go
```
