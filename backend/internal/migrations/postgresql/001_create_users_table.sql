CREATE TABLE IF NOT EXISTS Users (
    tg_id BIGINT UNIQUE NOT NULL,
    tg_tag VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT(CURRENT_TIMESTAMP) NOT NULL,
    last_messaged_at TIMESTAMP DEFAULT(CURRENT_TIMESTAMP) NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'requested',
	CONSTRAINT status_constraint CHECK(status IN('requested', 'confirmed', 'restricted', 'banned', 'admin'))
);

ALTER TABLE Users ADD COLUMN IF NOT EXISTS context VARCHAR;
ALTER TABLE Users ADD COLUMN IF NOT EXISTS context_data JSON;

ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (tg_id);
ALTER TABLE users DROP COLUMN id;