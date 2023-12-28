-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id, to_account_id ,amount
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListAccountTransactions :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY created_at DESC;

-- name: ListAccountCredits :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY created_at DESC;

-- name: ListTransactionFromOneAccountToOther :many
SELECT * FROM transfers
WHERE to_account_id = $1 AND from_account_id = $2
ORDER BY created_at DESC;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;



