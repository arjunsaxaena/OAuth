-- name: FindOrCreateUser :one
INSERT INTO users (name, phone, email, provider, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (phone , email)
DO UPDATE SET 
  phone = COALESCE(NULLIF(EXCLUDED.phone, ''), users.phone),
  provider = EXCLUDED.provider,
  name = COALESCE(NULLIF(EXCLUDED.name, ''), users.name),
  updated_at = NOW()
RETURNING id, phone, email, name, provider;
