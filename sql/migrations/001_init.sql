-- +goose Up
-- +goose StatementBegin
CREATE TABLE "bot_setup" (
                             "channel_id" text NOT NULL,
                             "message_id" text NOT NULL,
                             "created_at" timestamp NOT NULL DEFAULT now()
);
CREATE TABLE "player" (
                          "id" serial NOT NULL,
                          "name" character varying(255) NOT NULL,
                          PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "player_name_key" ON "player" ("name");
CREATE TABLE "char" (
                        "id" serial NOT NULL,
                        "name" character varying(255) NOT NULL,
                        "class" character varying(255) NOT NULL,
                        "user_id" integer NULL,
                        PRIMARY KEY ("id"),
                        CONSTRAINT "char_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE TABLE "debt" (
                        "id" serial NOT NULL,
                        "amount" bigint NOT NULL,
                        "last_updated" timestamp NOT NULL DEFAULT now(),
                        "user_id" integer NULL,
                        PRIMARY KEY ("id"),
                        CONSTRAINT "debt_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE UNIQUE INDEX "debt_user_id_key" ON "debt" ("user_id");
CREATE TABLE "debt_journal" (
                                "id" serial NOT NULL,
                                "amount" bigint NOT NULL,
                                "description" text NOT NULL DEFAULT '',
                                "date" timestamp NOT NULL DEFAULT now(),
                                "user_id" integer NULL,
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

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER "check_number_of_journal_rows" ON "debt_journal";
DROP FUNCTION "check_number_of_journal_rows"();
DROP TABLE "bot_setup";
DROP TABLE "player";
DROP TABLE "char";
DROP TABLE "debt";
DROP TABLE "debt_journal";
-- +goose StatementEnd
