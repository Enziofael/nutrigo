-- ============================================================
-- TABLE EXERCISES
-- ============================================================

CREATE TABLE IF NOT EXISTS training.exercises (
    "id" SERIAL NOT NULL CHECK("id" >= 0) PRIMARY KEY, -- uneditable (sets on creation only)
    "author_user_id" INT NOT NULL REFERENCES users("tg_id") ON DELETE CASCADE, -- uneditable (sets on creation only)
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- uneditable (sets on creation only)

    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- directly-uneditable (updates by trigger)
    "version" BIGINT NOT NULL DEFAULT 0 CHECK("version" >= 0), -- directly-uneditable (updates by trigger)
    "deleted_at" TIMESTAMPTZ DEFAULT NULL, -- directly-uneditable (sets by trigger)
    "last_entry_at" TIMESTAMPTZ DEFAULT NULL, -- directly-uneditable (updates by trigger)
    
    "personal_rating" INT NOT NULL DEFAULT 100 CHECK("rating" BETWEEN 0 AND 100), -- directly-uneditable (api-service AND trigger)

    "name" VARCHAR NOT NULL CHECK("name" >= 3), -- EDITABLE RACE-SIGNIFICANT
    "description" VARCHAR NOT NULL DEFAULT '', -- EDITABLE RACE-SIGNIFICANT
    "tags" VARCHAR[] NOT NULL DEFAULT '{}' -- EDITABLE RACE-SIGNIFICANT,
    "measurement_units" VARCHAR[] NOT NULL REFERENCES reference.measurement_units("code") -- EDITABLE RACE-SIGNIFICANT
);

COMMENT ON COLUMN training.exercises."id" IS 'uneditable (creation)';
COMMENT ON COLUMN training.exercises."user_id" IS 'uneditable (creation)';
COMMENT ON COLUMN training.exercises."created_at" IS 'uneditable (creation)';

COMMENT ON COLUMN training.exercises."updated_at" IS 'directly-uneditable (trigger)';
COMMENT ON COLUMN training.exercises."version" IS 'directly-uneditable (trigger)';
COMMENT ON COLUMN training.exercises."deleted_at" IS 'directly-uneditable (trigger)';
COMMENT ON COLUMN training.exercises."last_entry_at" IS 'directly-uneditable (trigger)';

COMMENT ON COLUMN training.exercises."personal_rating" IS 'directly-uneditable (api-service;trigger)';

COMMENT ON COLUMN training.exercises."name" IS 'editable (api)';
COMMENT ON COLUMN training.exercises."description" IS 'editable (api)';
COMMENT ON COLUMN training.exercises."tags" IS 'editable (api)';
COMMENT ON COLUMN training.exercises."measurement_units" IS 'editable (api)';

-- ============================================================
-- INDEXES
-- ============================================================

CREATE INDEX IF NOT EXISTS idx_exercises_user_id ON training.exercises(user_id);
CREATE INDEX IF NOT EXISTS idx_exercises_last_entry_at ON training.exercises(last_entry_at);
CREATE INDEX IF NOT EXISTS idx_exercises_personal_rating ON training.exercises(personal_rating);
-- search indexes
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises USING GIN (tags);
CREATE INDEX IF NOT EXISTS idx_exercises_name_trgm ON exercises USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_exercises_description_trgm ON exercises USING GIN (description gin_trgm_ops);

----------------------------------------------------------------------------------------
-- TABLE EXERCISE_ENTRIES

CREATE TABLE IF NOT EXISTS exercise_entries (
    id SERIAL PRIMARY KEY,
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    comment VARCHAR
);
-- Index by exercise_id
CREATE INDEX IF NOT EXISTS idx_exercise_entries_tg_id ON exercise_entries(exercise_id);

----------------------------------------------------------------------------------------
-- TABLE EXERCISE_SETS

CREATE TABLE IF NOT EXISTS exercise_sets (
    id SERIAL PRIMARY KEY,
    entry_id INT NOT NULL REFERENCES exercise_entries(id) ON DELETE CASCADE,
    measurements JSON NOT NULL,
    reps INT NOT NULL
);
-- Index by entry_id
CREATE INDEX IF NOT EXISTS idx_exercise_sets_tg_id ON exercise_sets(entry_id);







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


