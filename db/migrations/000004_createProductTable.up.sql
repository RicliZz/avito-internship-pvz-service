CREATE TYPE "type" AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE IF NOT EXISTS products (
    "ID" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "dateTime" timestamp NOT NULL DEFAULT now(),
    "type" type NOT NULL,
    "receptionID" uuid NOT NULL,
    FOREIGN KEY ("receptionID") REFERENCES reception ("ID") ON DELETE CASCADE
)