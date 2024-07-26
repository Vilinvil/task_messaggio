DROP TRIGGER IF EXISTS verify_updated_at_message ON public."message";
DROP FUNCTION IF EXISTS updated_at_now;

DROP TABLE IF EXISTS public."message";

DROP TYPE IF EXISTS public."message_status";

DROP TYPE IF EXISTS message_status;
