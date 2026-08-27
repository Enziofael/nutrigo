CREATE SCHEMA IF NOT EXISTS auth;

-- ============================================================
-- TABLE USERS
-- ============================================================

CREATE TABLE IF NOT EXISTS auth.users (
    "tg_id" BIGINT NOT NULL CHECK("tg_id" >= 0) PRIMARY KEY, -- directly-uneditable (sets on creation only)
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- directly-uneditable (sets on creation only)

    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- directly-uneditable (updates by trigger)
    "version" BIGINT NOT NULL DEFAULT 0 CHECK("version" >= 0), -- directly-uneditable (updates by trigger)
    "deleted_at" TIMESTAMPTZ DEFAULT NULL, -- directly-uneditable (sets by trigger)
    
    "last_messaged_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- directly-uneditable (updates by api-service by flag) race-insignificant

    "tg_tag" VARCHAR UNIQUE NOT NULL, -- EDITABLE RACE-SIGNIFICANT
    "account_status" VARCHAR NOT NULL DEFAULT 'requested' CHECK("account_status" IN('requested', 'confirmed', 'restricted', 'banned', 'admin')), -- EDITABLE RACE-SIGNIFICANT
    "tg_context_header" VARCHAR NOT NULL DEFAULT '', -- EDITABLE RACE-SIGNIFICANT
    "tg_context_data" JSON NOT NULL DEFAULT '{}', -- EDITABLE RACE-SIGNIFICANT
    "tg_context_message_id" BIGINT DEFAULT NULL -- EDITABLE NULLABLE RACE-SIGNIFICANT
);

INSERT INTO auth.users
(tg_id, tg_tag, account_status) VALUES
(0,'alanProg_bot','admin'); -- pre-made user for bot

-- ============================================================
-- INDEXES
-- ============================================================

CREATE INDEX IF NOT EXISTS idx_users_last_messaged_at ON auth.users (last_messaged_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at_null ON auth.users (deleted_at) WHERE deleted_at IS NULL;

-- ============================================================
-- DML TRIGGERS
-- ============================================================

-- ON INSERT
-- > protects default values
CREATE OR REPLACE FUNCTION auth.users_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- overwrite default values protection
    NEW.created_at = NOW();
    NEW.updated_at = NOW();
    NEW.version = 0;
    NEW.deleted_at = NULL;
    NEW.last_messaged_at = NOW();
    NEW.account_status = 'requested';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER auth.trg_users_on_insert
BEFORE INSERT ON auth.users
FOR EACH ROW
EXECUTE FUNCTION auth.users_on_insert();

-- ON UPDATE
-- > protects uneditable fields
-- > validates new values
-- > do all the logic about version, updated_at
CREATE OR REPLACE FUNCTION auth.users_on_update()
RETURNS TRIGGER AS $$
BEGIN
    -- change protection
    IF NEW.tg_id != OLD.tg_id THEN -- uneditable
        RAISE WARNING 'tg_id change attempted: % → % (ignored)', OLD.tg_id, NEW.tg_id;
        NEW.tg_id = OLD.tg_id;
    END IF;
    IF NEW.created_at != OLD.created_at THEN -- uneditable
        RAISE WARNING 'created_at change attempted: % → % (ignored)', OLD.created_at, NEW.created_at;
        NEW.created_at = OLD.created_at;
    END IF;
    IF NEW.version != OLD.version THEN -- uneditable (trigger)
        RAISE WARNING 'version change attempted: % → % (ignored)', OLD.version, NEW.version;
        NEW.version = OLD.version;
    END IF;
    IF NEW.updated_at != OLD.updated_at THEN -- uneditable (trigger)
        RAISE WARNING 'updated_at change attempted: % → % (ignored)', OLD.updated_at, NEW.updated_at;
        NEW.updated_at = OLD.updated_at;
    END IF;
    IF NEW.last_messaged_at < OLD.last_messaged_at THEN -- increase only
        RAISE WARNING 'last_messaged_at decrease attmepted:  % → % (ignored)', OLD.last_messaged_at, NEW.last_messaged_at;
        NEW.last_messaged_at = OLD.last_messaged_at;
    END IF;

    -- updated_at and version logic
    IF (OLD.tg_tag IS DISTINCT FROM NEW.tg_tag) OR
       (OLD.account_status IS DISTINCT FROM NEW.account_status) OR
       (OLD.tg_context_header IS DISTINCT FROM NEW.tg_context_header) OR
       (OLD.tg_context_data IS DISTINCT FROM NEW.tg_context_data) OR
       (OLD.tg_context_message_id IS DISTINCT FROM NEW.tg_context_message_id)
    THEN
        NEW.version = OLD.version + 1;
        NEW.updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER auth.trg_users_on_update
BEFORE UPDATE ON auth.users
FOR EACH ROW
EXECUTE FUNCTION auth.users_on_update();

-- ON DELETE
-- > soft deletes when deleted_at is NULL
-- > deletes permanently otherwise
CREATE OR REPLACE FUNCTION auth.users_on_delete()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.deleted_at IS NULL THEN
        UPDATE auth.users SET deleted_at = NOW() WHERE tg_id = OLD.tg_id;
        RAISE NOTICE 'user @% (id=%) soft deleted', OLD.tg_tag, OLD.tg_id;
        RETURN NULL;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER auth.trg_users_on_delete
BEFORE DELETE ON auth.users
FOR EACH ROW
EXECUTE FUNCTION auth.users_on_delete();





COMMENT ON COLUMN auth.users."tg_id" IS 'directly-uneditable (creation)';
COMMENT ON COLUMN auth.users."created_at" IS 'directly-uneditable (creation)';

COMMENT ON COLUMN auth.users."updated_at" IS 'directly-uneditable (trigger)';
COMMENT ON COLUMN auth.users."version" IS 'directly-uneditable (trigger)';
COMMENT ON COLUMN auth.users."deleted_at" IS 'directly-uneditable (trigger)';

COMMENT ON COLUMN auth.users."last_messaged_at" IS 'directly-uneditable (api-service)';

COMMENT ON COLUMN auth.users."tg_tag" IS 'editable (api)';
COMMENT ON COLUMN auth.users."account_status" IS 'editable (api)';
COMMENT ON COLUMN auth.users."tg_context_header" IS 'editable (api)';
COMMENT ON COLUMN auth.users."tg_context_data" IS 'editable (api)';
COMMENT ON COLUMN auth.users."tg_context_message_id" IS 'editable (api)';