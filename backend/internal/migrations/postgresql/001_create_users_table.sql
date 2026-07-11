CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE NOT NULL,
    tg_tag VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT(CURRENT_TIMESTAMP) NOT NULL,
    last_messaged_at TIMESTAMP DEFAULT(CURRENT_TIMESTAMP) NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'requested',
	CONSTRAINT status_constraint CHECK(status IN('requested', 'confirmed', 'blocked', 'banned', 'admin'))
);