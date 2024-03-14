CREATE TABLE player
(
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE char
(
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    class varchar(255) NOT NULL,
    user_id integer references player(id) ON DELETE CASCADE
);

CREATE TABLE debt
(
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    last_updated TIMESTAMP NOT NULL DEFAULT now(),
    user_id INTEGER UNIQUE REFERENCES player(id) ON DELETE CASCADE
);

CREATE TABLE debt_journal
(
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    date TIMESTAMP NOT NULL DEFAULT now(),
    user_id INTEGER REFERENCES player(id) ON DELETE CASCADE
);

CREATE TABLE bot_setup
(
    channel_id TEXT NOT NULL,
    message_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION check_number_of_journal_rows()
RETURNS TRIGGER AS 
$$
BEGIN
    IF (SELECT COUNT(*) FROM debt_journal WHERE user_id = NEW.user_id) > 9 THEN
        DELETE FROM debt_journal WHERE user_id in (SELECT user_id FROM debt_journal WHERE user_id = NEW.user_id ORDER BY date ASC LIMIT 1);
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER check_number_of_journal_rows
BEFORE INSERT ON debt_journal 
FOR EACH ROW EXECUTE FUNCTION check_number_of_journal_rows()
