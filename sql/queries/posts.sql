-- name: CreatePost :one
INSERT INTO posts (
	id,
	created_at,
	updated_at,
	title,
	url,
	description,
	published_at,
	feed_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning *;

-- name: GetPostsForUser :many
with feeds_for_user as (select * from feed_follows where feed_follows.user_id = $1)
SELECT posts.* FROM posts
join feeds_for_user on posts.feed_id = feeds_for_user.feed_id
order by published_at desc
limit $2;
