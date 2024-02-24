CREATE TABLE players
(
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE chars
(
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    class varchar(255) NOT NULL,
    user_id integer references players(id) ON DELETE CASCADE
);

CREATE TABLE debts
(
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    date DATE NOT NULL DEFAULT now(),
    user_id INTEGER REFERENCES players(id) ON DELETE CASCADE
);

