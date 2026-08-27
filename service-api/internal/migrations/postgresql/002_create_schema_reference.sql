CREATE SCHEMA IF NOT EXISTS reference;

--=====================================================
-- MEASUREMENT TYPES
--=====================================================

CREATE TABLE IF NOT EXISTS reference.measurement_types (
    "code" VARCHAR NOT NULL UNIQUE,

    "id" SERIAL NOT NULL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "version" BIGINT NOT NULL DEFAULT 0 CHECK("version" >= 0),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO reference.measurement_types ("code")
VALUES ('distance'), ('weight'), ('time'), ('speed') 
ON CONFLICT ("code") DO NOTHING;

-- ON INSERT
-- >> protects default values
CREATE OR REPLACE FUNCTION reference.measurement_types_on_insert ()
RETURNS TRIGGER AS $$
    -- overwrite default values protection
    NEW."created_at" = NOW();
    NEW."updated_at" = NOW();
    NEW."version" = 0;
    NEW."deleted_at" = NULL;
    RETURN NEW;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_types_on_insert
BEFORE INSERT ON reference.measurement_types
FOR EACH ROW EXECUTE FUNCTION reference.measurement_types_on_insert();

-- ON UPDATE
-- >> protects uneditable fields
-- >> do all the logic about version, updated_at
CREATE OR REPLACE FUNCTION reference.measurement_types_on_update ()
RETURNS TRIGGER AS $$
    IF NEW."id" != OLD."id" THEN -- uneditable
        RAISE WARNING 'measurement_type id change attempted: % → % (ignored)', OLD."id", NEW."id";
        NEW."id" = OLD."id";
    END IF;
    IF NEW."created_at" != OLD."created_at" THEN -- uneditable
        RAISE WARNING 'measurement_type created_at change attempted: % → % (ignored)', OLD."created_at", NEW."created_at";
        NEW."created_at" = OLD."created_at";
    END IF;
    IF NEW."updated_at" != OLD."updated_at" THEN -- uneditable
        RAISE WARNING 'measurement_type updated_at change attempted: % → % (ignored)', OLD."updated_at", NEW."updated_at";
        NEW."updated_at" = OLD."updated_at";
    END IF;
    IF NEW."version" != OLD."version" THEN -- uneditable
        RAISE WARNING 'measurement_type version change attempted: % → % (ignored)', OLD."version", NEW."version";
        NEW."version" = OLD."version";
    END IF;
    IF (OLD."deleted_at" IS NOT NULL) AND (NEW."deleted_at" != OLD."deleted_at") THEN -- uneditable
        RAISE WARNING 'measurement_type deleted_at change attempted: % → % (ignored)', OLD."deleted_at", NEW."deleted_at";
        NEW."deleted_at" = OLD."deleted_at";
    END IF;
    IF (OLD."code" IS DISTINCT FROM NEW."code") THEN
        NEW."version" = OLD."version" + 1;
        NEW."updated_at" = NOW();
    END IF;
    RETURN NEW;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_types_on_update
BEFORE UPDATE ON reference.measurement_types
FOR EACH ROW EXECUTE FUNCTION reference.measurement_types_on_update();

-- ON DELETE
-- >> soft deletes when deleted_at is NULL
-- >> deletes permanently otherwise
CREATE OR REPLACE FUNCTION reference.measurement_types_on_delete ()
RETURNS TRIGGER AS $$
    IF OLD."deleted_at" IS NULL THEN
        UPDATE reference.measurement_types SET "deleted_at" = NOW() WHERE "id" = OLD."id";
        RAISE NOTICE 'measurement_type "%" (id=%) soft deleted', OLD."code", OLD."id";
        RETURN NULL;
    ELSE
        RETURN OLD;
    END IF;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_types_on_delete
BEFORE DELETE ON reference.measurement_types
FOR EACH ROW EXECUTE FUNCTION reference.measurement_types_on_delete();

-- SELECT * FROM reference.measurement_types;

--=====================================================
-- MEASUREMENT UNITS
--=====================================================

CREATE TABLE IF NOT EXISTS reference.measurement_units (
    "code" VARCHAR NOT NULL UNIQUE,
    "type_code" VARCHAR NOT NULL REFERENCES reference.measurement_types("code") ON UPDATE CASCADE ON DELETE RESTRICT,
    "base_unit" VARCHAR NOT NULL REFERENCES reference.measurement_units("code") ON UPDATE CASCADE,
    "unit_to_base_ratio" NUMERIC NOT NULL,

    "id" SERIAL NOT NULL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "version" BIGINT NOT NULL DEFAULT 0 CHECK("version" >= 0),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO reference.measurement_units ("type_code", "code", "base_unit", "unit_to_base_ratio")
VALUES 
('distance', 'km', 'km', 1.0), ('distance', 'mi', 'km', 1.609344),
('weight', 'kg', 'kg', 1.0), ('weight', 'lb', 'kg', 0.45359237),
('time', 's', 's', 1.0), ('time', 'min', 's', 60.0), ('time', 'h', 's', 3600.0), ('time', 'd', 's', 86400.0),  
('speed', 'm_s', 'm_s', 1.0), ('speed', 'km_h', 'm_s', 0.2777777778), ('speed', 'mi_h', 'm_s', 0.44704), ('speed', 'kn', 'm_s', 0.5144444444)
ON CONFLICT ("code") DO NOTHING;

-- ON INSERT
-- >> protects default values
CREATE OR REPLACE FUNCTION reference.measurement_units_on_insert ()
RETURNS TRIGGER AS $$
    -- overwrite default values protection
    NEW."created_at" = NOW();
    NEW."updated_at" = NOW();
    NEW."version" = 0;
    NEW."deleted_at" = NULL;
    RETURN NEW;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_units_on_insert
BEFORE INSERT ON reference.measurement_units
FOR EACH ROW EXECUTE FUNCTION reference.measurement_units_on_insert();

-- ON UPDATE
-- >> protects uneditable fields
-- >> do all the logic about version, updated_at
CREATE OR REPLACE FUNCTION reference.measurement_units_on_update ()
RETURNS TRIGGER AS $$
    IF NEW."id" != OLD."id" THEN -- uneditable
        RAISE WARNING 'measurement_unit id change attempted: % → % (ignored)', OLD."id", NEW."id";
        NEW."id" = OLD."id";
    END IF;
    IF NEW."created_at" != OLD."created_at" THEN -- uneditable
        RAISE WARNING 'measurement_unit created_at change attempted: % → % (ignored)', OLD."created_at", NEW."created_at";
        NEW."created_at" = OLD."created_at";
    END IF;
    IF NEW."updated_at" != OLD."updated_at" THEN -- uneditable
        RAISE WARNING 'measurement_unit updated_at change attempted: % → % (ignored)', OLD."updated_at", NEW."updated_at";
        NEW."updated_at" = OLD."updated_at";
    END IF;
    IF NEW."version" != OLD."version" THEN -- uneditable
        RAISE WARNING 'measurement_unit version change attempted: % → % (ignored)', OLD."version", NEW."version";
        NEW."version" = OLD."version";
    END IF;
    IF (OLD."deleted_at" IS NOT NULL) AND (NEW."deleted_at" != OLD."deleted_at") THEN -- uneditable
        RAISE WARNING 'measurement_unit deleted_at change attempted: % → % (ignored)', OLD."deleted_at", NEW."deleted_at";
        NEW."deleted_at" = OLD."deleted_at";
    END IF;
    IF (OLD."code" IS DISTINCT FROM NEW."code") OR
       (OLD."type_code" IS DISTINCT FROM NEW."type_code") OR
       (OLD."base_unit" IS DISTINCT FROM NEW."base_unit") OR
       (OLD."unit_to_base_ratio" IS DISTINCT FROM NEW."unit_to_base_ratio") 
    THEN
        NEW."version" = OLD."version" + 1;
        NEW."updated_at" = NOW();
    END IF;
    RETURN NEW;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_units_on_update
BEFORE UPDATE ON reference.measurement_units
FOR EACH ROW EXECUTE FUNCTION reference.measurement_units_on_update();

-- ON DELETE
-- >> soft deletes when deleted_at is NULL
-- >> deletes permanently otherwise
CREATE OR REPLACE FUNCTION reference.measurement_units_on_delete ()
RETURNS TRIGGER AS $$
    IF OLD."deleted_at" IS NULL THEN
        UPDATE reference.measurement_units SET "deleted_at" = NOW() WHERE "id" = OLD."id";
        RAISE NOTICE 'measurement_unit "%" (id=%) soft deleted', OLD."code", OLD."id";
        RETURN NULL;
    ELSE
        RETURN OLD;
    END IF;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_units_on_delete
BEFORE DELETE ON reference.measurement_units
FOR EACH ROW EXECUTE FUNCTION reference.measurement_units_on_delete();

-- SELECT * FROM reference.measurement_units;

--=====================================================
-- LOCALES
--=====================================================
--BCP 47
CREATE TABLE IF NOT EXISTS reference.locales ( 
    "language" VARCHAR(8) NOT NULL,
    "script" VARCHAR(4),
    "region" VARCHAR(3),
    "variant" VARCHAR(8),
    UNIQUE ("language", "script", "region", "variant"),
    "full_tag" VARCHAR(35) PRIMARY KEY GENERATED ALWAYS AS (
        "language" || 
        CASE WHEN "script" IS NOT NULL THEN '-' || "script" ELSE '' END ||
        CASE WHEN "region" IS NOT NULL THEN '-' || "region" ELSE '' END ||
        CASE WHEN "variant" IS NOT NULL THEN '-' || "variant" ELSE '' END) STORED,
    "native_name" VARCHAR(50) NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT FALSE,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "version" BIGINT NOT NULL DEFAULT 0 CHECK("version" >= 0),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO reference.locales ("language")
VALUES ('ru', 'en')
ON CONFLICT ("language") DO NOTHING;

-- ON INSERT
-- >> protects default values
CREATE OR REPLACE FUNCTION reference.measurement_units_on_insert ()
RETURNS TRIGGER AS $$
    -- overwrite default values protection
    NEW."created_at" = NOW();
    NEW."updated_at" = NOW();
    NEW."version" = 0;
    NEW."deleted_at" = NULL;
    RETURN NEW;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_units_on_insert
BEFORE INSERT ON reference.measurement_units
FOR EACH ROW EXECUTE FUNCTION reference.measurement_units_on_insert();

-- ON UPDATE
-- >> protects uneditable fields
-- >> do all the logic about version, updated_at
CREATE OR REPLACE FUNCTION reference.measurement_units_on_update ()
RETURNS TRIGGER AS $$
    IF NEW."id" != OLD."id" THEN -- uneditable
        RAISE WARNING 'measurement_unit id change attempted: % → % (ignored)', OLD."id", NEW."id";
        NEW."id" = OLD."id";
    END IF;
    IF NEW."created_at" != OLD."created_at" THEN -- uneditable
        RAISE WARNING 'measurement_unit created_at change attempted: % → % (ignored)', OLD."created_at", NEW."created_at";
        NEW."created_at" = OLD."created_at";
    END IF;
    IF NEW."updated_at" != OLD."updated_at" THEN -- uneditable
        RAISE WARNING 'measurement_unit updated_at change attempted: % → % (ignored)', OLD."updated_at", NEW."updated_at";
        NEW."updated_at" = OLD."updated_at";
    END IF;
    IF NEW."version" != OLD."version" THEN -- uneditable
        RAISE WARNING 'measurement_unit version change attempted: % → % (ignored)', OLD."version", NEW."version";
        NEW."version" = OLD."version";
    END IF;
    IF (OLD."deleted_at" IS NOT NULL) AND (NEW."deleted_at" != OLD."deleted_at") THEN -- uneditable
        RAISE WARNING 'measurement_unit deleted_at change attempted: % → % (ignored)', OLD."deleted_at", NEW."deleted_at";
        NEW."deleted_at" = OLD."deleted_at";
    END IF;
    IF (OLD."code" IS DISTINCT FROM NEW."code") OR
       (OLD."type_code" IS DISTINCT FROM NEW."type_code") OR
       (OLD."base_unit" IS DISTINCT FROM NEW."base_unit") OR
       (OLD."unit_to_base_ratio" IS DISTINCT FROM NEW."unit_to_base_ratio") 
    THEN
        NEW."version" = OLD."version" + 1;
        NEW."updated_at" = NOW();
    END IF;
    RETURN NEW;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_units_on_update
BEFORE UPDATE ON reference.measurement_units
FOR EACH ROW EXECUTE FUNCTION reference.measurement_units_on_update();

-- ON DELETE
-- >> soft deletes when deleted_at is NULL
-- >> deletes permanently otherwise
CREATE OR REPLACE FUNCTION reference.measurement_units_on_delete ()
RETURNS TRIGGER AS $$
    IF OLD."deleted_at" IS NULL THEN
        UPDATE reference.measurement_units SET "deleted_at" = NOW() WHERE "id" = OLD."id";
        RAISE NOTICE 'measurement_unit "%" (id=%) soft deleted', OLD."code", OLD."id";
        RETURN NULL;
    ELSE
        RETURN OLD;
    END IF;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trg_measurement_units_on_delete
BEFORE DELETE ON reference.measurement_units
FOR EACH ROW EXECUTE FUNCTION reference.measurement_units_on_delete();

--=====================================================
-- LOCALIZATION FORM
--=====================================================

CREATE TABLE IF NOT EXISTS reference.localization_forms (
    "form" VARCHAR NOT NULL PRIMARY KEY
);

INSERT INTO reference.localization_forms ("form")
VALUES ('default'), ('short'), ('singular'), ('few'), ('many'), ('other')
ON CONFLICT DO NOTHING;

--=====================================================
-- LOCALIZATION ENTITIES
--=====================================================

CREATE TABLE IF NOT EXISTS reference.localization_entity_types (
    "type" VARCHAR NOT NULL PRIMARY KEY
);

INSERT INTO reference.localization_entities ("type")
VALUES ('measurement_type'), ('measurement_unit')
ON CONFLICT DO NOTHING;

--=====================================================
-- LOCALIZATIONS
--=====================================================

CREATE TABLE IF NOT EXISTS reference.localizations (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "entity_type" VARCHAR NOT NULL REFERENCES reference.localization_entities("type") ON UPDATE CASCADE ON DELETE RESTRICT,
    "entity_id" INT NOT NULL REFERENCES reference.measurement_units("code") ON UPDATE CASCADE ON DELETE RESTRICT,
    "locale" VARCHAR(10) NOT NULL REFERENCES reference.locales("full_tag") ON UPDATE CASCADE ON DELETE RESTRICT,
    "form" VARCHAR NOT NULL REFERENCES reference.localization_forms("form") ON UPDATE CASCADE ON DELETE RESTRICT,
    UNIQUE("entity_type","entity_id","locale","form"),
    "localization" VARCHAR NOT NULL,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "version" BIGINT NOT NULL DEFAULT 0 CHECK("version" >= 0),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL
);

INSERT INTO reference.localizations
("entity_type", "entity_id", "locale", "form", "localization")
VALUES
('measurement_type', 1, 'ru', 'default', 'дистанция'),
('measurement_type', 1, 'en', 'default', 'distance'),

('measurement_type', 2, 'ru', 'default', 'вес'),
('measurement_type', 2, 'en', 'default', 'weight'),

('measurement_type', 3, 'ru', 'default', 'время'), 
('measurement_type', 3, 'en', 'default', 'time'),

('measurement_type', 4, 'ru', 'default', 'скорость'),
('measurement_type', 4, 'en', 'default', 'speed')
ON CONFLICT("entity_type", "entity_id", "locale", "form", "localization") DO NOTHING;

USING (
SELECT

)
INSERT INTO reference.localizations 
("entity_type", "entity_id", "locale", "form", "localization")
VALUES 
('measurement_unit', 'km',       'ru', 'short', 'short', 'км'),
('measurement_unit', 'mi',       'ru', 'short', 'миля'),
('measurement_unit', 'kg',       'ru', 'short', 'кг'),
('measurement_unit', 'lb',       'ru', 'short', 'фунт'),
('measurement_unit', 's',        'ru', 'short', 'с'),
('measurement_unit', 'min',      'ru', 'short', 'мин'),
('measurement_unit', 'h',        'ru', 'short', 'ч'),
('measurement_unit', 'm_s',      'ru', 'short', 'м/с'),
('measurement_unit', 'km_h',     'ru', 'short', 'км/ч'),
('measurement_unit', 'mi_h',     'ru', 'short','миля/ч'),
('measurement_unit', 'km',       'en', 'short', 'km'),
('measurement_unit', 'mi',       'en', 'short', 'mi'),
('measurement_unit', 'kg',       'en', 'short','kg'),
('measurement_unit', 'lb',       'en', 'short', 'lb'),
('measurement_unit', 's',        'en', 'short', 's'),
('measurement_unit', 'min',      'en', 'short','min'),
('measurement_unit', 'h',        'en', 'short', 'h'),
('measurement_unit', 'm_s',      'en', 'short', 'm/s'),
('measurement_unit', 'km_h',     'en', 'short', 'km/h'),
('measurement_unit', 'mi_h',     'en', 'short', 'mph')
ON CONFLICT ("code", "locale", "form") DO NOTHING;
