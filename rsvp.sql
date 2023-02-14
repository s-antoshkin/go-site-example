DROP TABLE IF EXISTS rsvp;

CREATE TABLE IF NOT EXISTS rsvp (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    will_attend BOOLEAN NOT NULL
);