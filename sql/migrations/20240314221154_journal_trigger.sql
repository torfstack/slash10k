-- Modify "debt_journal" table
ALTER TABLE "debt_journal" ALTER COLUMN "date" TYPE timestamp;

-- Create "check_number_of_journal_rows" function
CREATE FUNCTION "check_number_of_journal_rows" () RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF (SELECT COUNT(*) FROM debt_journal WHERE user_id = NEW.user_id) > 9 THEN
        DELETE FROM debt_journal WHERE (user_id, date) in
            (SELECT user_id, min(date) FROM debt_journal WHERE user_id = NEW.user_id GROUP BY user_id);
    END IF;
    RETURN NEW;
END;
$$;
-- Create trigger "check_number_of_journal_rows"
CREATE TRIGGER "check_number_of_journal_rows" BEFORE INSERT ON "debt_journal" FOR EACH ROW EXECUTE FUNCTION "check_number_of_journal_rows"();
