-- +goose Up
ALTER TABLE feeds
	 add column last_fetched_at timestamp;

-- +goose Down
ALTER TABLE feeds
	 drop column last_fetched_at;
