CREATE TABLE IF NOT EXISTS chars(
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    class varchar(255) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    user_id integer references players(id) ON DELETE CASCADE
);

CREATE SEQUENCE IF NOT EXISTS chars_mig_seq START WITH 1 INCREMENT BY 1;
