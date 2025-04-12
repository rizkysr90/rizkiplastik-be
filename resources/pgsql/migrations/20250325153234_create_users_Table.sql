-- migrate:up
-- Create the users table
CREATE TYPE roles AS ENUM ('ADMIN', 'STAFF');

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(30) NOT NULL,
    hash_password TEXT NOT NULL,
    role roles NOT NULL,
    token TEXT,
    last_login TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (username)
);

-- migrate:down

