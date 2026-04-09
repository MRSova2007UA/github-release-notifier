CREATE TABLE IF NOT EXISTS subscribers (
                                           id SERIAL PRIMARY KEY,
                                           email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE TABLE IF NOT EXISTS repositories (
                                            id SERIAL PRIMARY KEY,
                                            name VARCHAR(255) UNIQUE NOT NULL,
    last_seen_tag VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE TABLE IF NOT EXISTS subscriptions (
                                             subscriber_id INT REFERENCES subscribers(id) ON DELETE CASCADE,
    repository_id INT REFERENCES repositories(id) ON DELETE CASCADE,
    PRIMARY KEY (subscriber_id, repository_id)
    );