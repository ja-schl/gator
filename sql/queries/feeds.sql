-- name: CreateFeed :one
INSERT INTO feeds (
	id, created_at, updated_at, name, url, user_id
) VALUES ( $1, $2, $3, $4, $5, $6) returning *;

-- name: GetFeeds :many
Select * from feeds;

-- name: CreateFeedFollow :one
with inserted_feed_follows as (INSERT INTO feed_follows (
	id, created_at, updated_at, feed_id, user_id
) VALUES ( $1, $2, $3, $4, $5 ) returning *)
select inserted_feed_follows.*, feeds.name as feed_name, users.name as user_name
	from inserted_feed_follows
inner join feeds on inserted_feed_follows.feed_id = feeds.id
inner join users on inserted_feed_follows.user_id = users.id;


-- name: GetFeedByUrl :one
SELECT * FROM feeds
where url = $1;

-- name: GetFeedFollowsForUser :many
with feeds_for_user as (select * from feed_follows where feed_follows.user_id = $1)
SELECT feeds_for_user.*, users.name as user_name, feeds.name as feed_name FROM feeds_for_user
inner join users on users.id = feeds_for_user.user_id
inner join feeds on feeds.id = feeds_for_user.feed_id;

-- name: DeleteFollow :exec
Delete from feed_follows where user_id = $1 and feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
	SET last_fetched_at = (select current_timestamp),
updated_at = (select current_timestamp)
	WHERE id = $1;

-- name: GetNextFeedToFetch :one
Select * from feeds 
order by last_fetched_at asc nulls first;
