-- name: CreateUser :one
INSERT INTO users (nickname, img, country, city)
VALUES ($1, $2, $3, $4)
RETURNING id, nickname, img, country, city, created_at, updated_at, deleted;

-- name: UpdateUser :one
UPDATE users
SET country = $2,
    city = $3,
    updated_at = now()
WHERE id = $1 AND deleted = FALSE
RETURNING id, nickname, img, country, city, created_at, updated_at, deleted;

-- name: UpdateUserImg :one
UPDATE users
SET img = $2,
    updated_at = now()
WHERE id = $1 AND deleted = FALSE
RETURNING id, nickname, img, country, city, created_at, updated_at, deleted;

-- name: SoftDeleteUser :one
UPDATE users
SET deleted = TRUE,
    updated_at = now()
WHERE id = $1
RETURNING id, nickname, img, country, city, created_at, updated_at, deleted;

-- name: GetUserByID :one
SELECT u.id, u.nickname, u.img, u.country, u.city, u.created_at, u.updated_at, u.deleted
FROM users u
WHERE u.id = $1;

-- name: GetUserImgByNickname :one
SELECT img
FROM users
WHERE nickname = $1 AND deleted = FALSE;

-- name: ListUsers :many
SELECT u.id, u.nickname, u.img, u.country, u.city, u.created_at, u.updated_at, u.deleted
FROM users u
WHERE u.deleted = FALSE
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: FindUsers :many
SELECT u.id, u.nickname, u.img, u.country, u.city, u.created_at, u.updated_at, u.deleted
FROM users u
WHERE u.deleted = FALSE
ORDER BY u.created_at DESC;

-- name: CreateClub :one
INSERT INTO clubs (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
RETURNING id, name;

-- name: AddUserClub :exec
INSERT INTO user_clubs (user_id, club_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteUserClubsByUserID :exec
DELETE FROM user_clubs
WHERE user_id = $1;

-- name: GetClubsByUserID :many
SELECT c.id, c.name
FROM clubs c
JOIN user_clubs uc ON uc.club_id = c.id
WHERE uc.user_id = $1;
