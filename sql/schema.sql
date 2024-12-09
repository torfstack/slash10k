CREATE TABLE player
(
    id SERIAL PRIMARY KEY,
    discord_id TEXT NOT NULL,
    discord_name TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    name varchar(255) NOT NULL,
    UNIQUE (discord_id, guild_id)
);

CREATE TABLE debt
(
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    last_updated TIMESTAMP NOT NULL DEFAULT now(),
    user_id INTEGER UNIQUE NOT NULL REFERENCES player(id) ON DELETE CASCADE
);

CREATE TABLE debt_journal
(
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    date TIMESTAMP NOT NULL DEFAULT now(),
    user_id INTEGER NOT NULL REFERENCES player(id) ON DELETE CASCADE
);

CREATE TABLE bot_setup
(
    guild_id TEXT UNIQUE NOT NULL,
    channel_id TEXT NOT NULL,
    registration_message_id TEXT NOT NULL,
    debts_message_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION check_number_of_journal_rows()
RETURNS TRIGGER AS 
$$
BEGIN
    IF (SELECT COUNT(*) FROM debt_journal WHERE user_id = NEW.user_id) > 9 THEN
        DELETE FROM debt_journal WHERE user_id in (SELECT user_id FROM debt_journal WHERE user_id = NEW.user_id ORDER BY date ASC LIMIT 1 AND guild_id = NEW.guild_id);
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER check_number_of_journal_rows
BEFORE INSERT ON debt_journal 
FOR EACH ROW EXECUTE FUNCTION check_number_of_journal_rows();

CREATE OR REPLACE FUNCTION create_debt_for_new_player()
RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO debt (amount, user_id)
    VALUES (0, NEW.id);
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER create_debt_for_new_player
AFTER INSERT ON player
FOR EACH ROW EXECUTE FUNCTION create_debt_for_new_player();