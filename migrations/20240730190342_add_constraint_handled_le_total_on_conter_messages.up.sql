ALTER TABLE public."counter_message" DROP CONSTRAINT IF EXISTS handled_le_total;
ALTER TABLE  public."counter_message" ADD CONSTRAINT handled_le_total CHECK (handled <= total);