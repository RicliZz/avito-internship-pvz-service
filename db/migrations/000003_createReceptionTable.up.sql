CREATE TYPE status AS ENUM('in_progress', 'close');

CREATE TABLE IF NOT EXISTS reception(
    "ID" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "dateTime" timestamp NOT NULL DEFAULT now(),
    "pvzID" uuid NOT NULL,
    status status NOT NULL DEFAULT 'in_progress',
    FOREIGN KEY ("pvzID") REFERENCES "PVZ" ("ID") ON DELETE CASCADE
)