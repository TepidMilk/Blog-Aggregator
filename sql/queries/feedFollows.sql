-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted_feed_follow.*,
feeds.name as feed_name,
users.name as user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT users.name as user_name, feeds.name as feed_name, feed_follows.* FROM feed_follows
INNER JOIN users on users.id = feed_follows.user_id
INNER JOIN feeds on feeds.id = feed_follows.feed_id
WHERE users.name = $1;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;