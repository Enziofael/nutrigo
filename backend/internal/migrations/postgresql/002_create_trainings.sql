CREATE TABLE IF NOT EXISTS exercises (
    id SERIAL PRIMARY KEY,
    tg_id INT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    technique VARCHAR,
    weight_unit VARCHAR NOT NULL DEFAULT 'kg'
    rating INT NOT NULL DEFAULT 100
);

ALTER TABLE exercises
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_exercises_tg_id ON exercises(tg_id);

CREATE TABLE IF NOT EXISTS exercise_entries (
    id SERIAL PRIMARY KEY,
    tg_id INT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    comment VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_exercise_entries_tg_id ON exercise_entries(tg_id);

CREATE TABLE IF NOT EXISTS exercise_sets (
    id SERIAL PRIMARY KEY,
    tg_id INT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    entry_id INT NOT NULL REFERENCES exercise_entries(id) ON DELETE CASCADE,
    weight_kg DECIMAL NOT NULL,
    weight_lbs DECIMAL NOT NULL,
    reps INT NOT NULL,
    weight_unit VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_exercise_sets_tg_id ON exercise_sets(tg_id);

CREATE TABLE IF NOT EXISTS training_programms (
    id SERIAL PRIMARY KEY,
    tg_id INT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    description VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_training_programms_tg_id ON training_programms(tg_id);

CREATE TABLE IF NOT EXISTS training_programm_days (
    id SERIAL PRIMARY KEY,
    tg_id INT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    training_programm_id INT NOT NULL REFERENCES training_programms(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    period INT,
    description VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_training_programm_days_tg_id ON training_programm_days(tg_id);

CREATE TABLE IF NOT EXISTS training_program_days_exercises (
    training_programm_day_id INT NOT NULL REFERENCES training_programm_days(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    PRIMARY KEY (training_programm_day_id, exercise_id)
);


-- для поиска
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_exercises_name_trgm ON exercises USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_exercises_description_trgm ON exercises USING GIN (description gin_trgm_ops);