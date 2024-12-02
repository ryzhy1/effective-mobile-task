-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    group_name   VARCHAR(255) NOT NULL,
    title        VARCHAR(255) NOT NULL,
    release_date DATE,
    text         TEXT,
    link         TEXT,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd
