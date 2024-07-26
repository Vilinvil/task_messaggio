DO $$
    BEGIN
        IF NOT EXISTS (SELECT * FROM pg_type WHERE typname = 'message_status') THEN
            create type message_status AS ENUM ('pending', 'done');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS public."message"
(
    id UUID PRIMARY KEY,
    value TEXT NOT NULL
        CONSTRAINT len_value_not_zero CHECK (value <> '')
        CONSTRAINT max_len_value CHECK (LENGTH(value) <= 1000),
    status public."message_status" NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION updated_at_now()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS verify_updated_at_message ON public."message";
CREATE TRIGGER verify_updated_at_message
    BEFORE UPDATE
    ON public."message"
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_now();
