CREATE TABLE IF NOT EXISTS counter (
                         id serial PRIMARY KEY,
                         name VARCHAR (128) UNIQUE NOT NULL,
                         value BIGINT NOT NULL
);
CREATE TABLE IF NOT EXISTS gauge (
                         id serial PRIMARY KEY,
                         name VARCHAR (128) UNIQUE NOT NULL,
                         value DOUBLE PRECISION NOT NULL
);