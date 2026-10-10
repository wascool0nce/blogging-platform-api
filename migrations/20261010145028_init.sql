-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    title varchar(100) NOT NULL,
    content TEXT,
    category varchar(100) NOT NULL,
    tags varchar(50)[],
    created_at timestamp NOT NULL,
    updated_at timestamp
);

-- +goose Down
DROP TABLE post;
