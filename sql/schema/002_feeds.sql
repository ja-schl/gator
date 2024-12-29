-- +goose Up
create table feeds (
id uuid PRIMARY KEY,
created_at timestamp not null,
updated_at timestamp not null,
name text not null,
url text unique not null, 
user_id uuid not null,
constraint fk_user_id
	FOREIGN KEY(user_id) REFERENCES users(id) on delete cascade
);

-- +goose Down
DROP TABLE feeds;
