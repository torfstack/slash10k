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
    date DATE NOT NULL DEFAULT now(),
    user_id INTEGER REFERENCES player(id) ON DELETE CASCADE
);
