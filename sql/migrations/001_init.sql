-- +goose Up
-- +goose StatementBegin
CREATE TABLE "bot_setup" (
                             "guild_id" text UNIQUE NOT NULL,
                             "channel_id" text NOT NULL,
                             "debts_message_id" text NOT NULL,
                             "registration_message_id" text NOT NULL,
                             "created_at" timestamp NOT NULL DEFAULT now()
);
CREATE TABLE "player" (
                          "id" serial NOT NULL,
                          "discord_id" text NOT NULL,
                          "discord_name" text NOT NULL,
                          "guild_id" text NOT NULL,
                          "name" character varying(255) NOT NULL,
                          PRIMARY KEY ("id"),
                          UNIQUE ("discord_id", "guild_id")
);
CREATE TABLE "debt" (
                        "id" serial NOT NULL,
                        "amount" bigint NOT NULL,
                        "last_updated" timestamp NOT NULL DEFAULT now(),
                        "user_id" integer UNIQUE NOT NULL,
                        PRIMARY KEY ("id"),
                        CONSTRAINT "debt_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
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

CREATE FUNCTION create_debt_for_new_player () RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO debt (amount, user_id)
    VALUES (0, NEW.id);
    RETURN NEW;
END;
$$;
CREATE TRIGGER create_debt_for_new_player AFTER INSERT ON player FOR EACH ROW EXECUTE FUNCTION create_debt_for_new_player();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER "create_debt_for_new_player" ON "player";
DROP FUNCTION "create_debt_for_new_player"();
DROP TRIGGER "check_number_of_journal_rows" ON "debt_journal";
DROP FUNCTION "check_number_of_journal_rows"();
DROP TABLE "bot_setup";
DROP TABLE "player";
DROP TABLE "debt";
DROP TABLE "debt_journal";
-- +goose StatementEnd
