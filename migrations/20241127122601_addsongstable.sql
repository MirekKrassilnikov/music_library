-- +goose Up
CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       group_name VARCHAR(255) NOT NULL,
                       song_name VARCHAR(255) NOT NULL,
                       release_date DATE,
                       text TEXT,
                       link TEXT
);

-- +goose Down
DROP TABLE IF EXISTS songs;
