-- +goose Up
-- +goose StatementBegin
DROP TRIGGER "check_number_of_journal_rows" ON "debt_journal";
DROP FUNCTION "check_number_of_journal_rows"();
DROP TABLE "debt_journal";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE "debt_journal" (
    "id" serial NOT NULL,
    "amount" bigint NOT NULL,
    "description" text NOT NULL DEFAULT '',
    "date" timestamp NOT NULL DEFAULT now(),
    "user_id" integer NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "debt_journal_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE FUNCTION "check_number_of_journal_rows" () RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF (SELECT COUNT(*) FROM debt_journal WHERE user_id = NEW.user_id) > 9 THEN
        DELETE FROM debt_journal WHERE (user_id, date) in
                                       (SELECT user_id, min(date) FROM debt_journal WHERE user_id = NEW.user_id GROUP BY user_id);
    END IF;
    RETURN NEW;
END;
$$;
CREATE TRIGGER "check_number_of_journal_rows" BEFORE INSERT ON "debt_journal" FOR EACH ROW EXECUTE FUNCTION "check_number_of_journal_rows"();
-- +goose StatementEnd
