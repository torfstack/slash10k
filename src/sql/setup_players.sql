CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL
);

CREATE SEQUENCE IF NOT EXISTS players_mig_seq START WITH 1 INCREMENT BY 1