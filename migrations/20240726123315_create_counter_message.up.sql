CREATE TABLE IF NOT EXISTS public."counter_message"
(
    id         BIGINT DEFAULT 0 NOT NULL PRIMARY KEY
        CONSTRAINT only_one_tuple CHECK (id = 0),
    total      BIGINT DEFAULT 0 NOT NULL
        CONSTRAINT not_negative_total CHECK (total >= 0),
    handled BIGINT DEFAULT 0 NOT NULL
        CONSTRAINT not_negative_handled CHECK (handled >= 0)
);

INSERT INTO public."counter_message" (id, total, handled) VALUES (0, 0, 0);
