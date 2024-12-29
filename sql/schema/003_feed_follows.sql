-- +goose Up
create table feed_follows(
id uuid primary key,
	created_at timestamp not null,
	updated_at timestamp not null,
	feed_id uuid not null references feeds(id)  on delete cascade,
	user_id uuid not null references users(id)  on delete cascade,
	unique(feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;
