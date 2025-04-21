-- name: StoreToken :one
INSERT INTO tokens (
    user_id,jti,expires_at
) VALUES ( ?,?,?)
RETURNING *     ;

-- name: IsRevoked :one
SELECT is_revoked FROM tokens
WHERE jti = ?;

-- name: RevokeToken :one
UPDATE tokens 
SET is_revoked  = 1
WHERE jti = ?
RETURNING * ;