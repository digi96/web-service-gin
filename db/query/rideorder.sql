-- name: CreateOrder :one
INSERT INTO rideorder(
    contact_id,
    rider_name,
    rider_phone,
    destination
  )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetOrderById :one
SELECT *
FROM rideorder
WHERE rideorder_id = $1
LIMIT 1;
-- name: ListOrders :many
SELECT *
FROM rideorder
ORDER BY rideorder_id
LIMIT $1 OFFSET $2;
-- name: UpdateOrder :one
UPDATE rideorder
SET pickup_at = coalesce(sqlc.narg('pickup_at'), pickup_at),
  updated_at = coalesce(sqlc.narg('updated_at'), updated_at)
WHERE rideorder_id = sqlc.arg('riderorder_id')
RETURNING *;
-- name: DeleteOrder :exec
DELETE FROM rideorder
WHERE rideorder_id = $1;