CREATE TYPE status AS ENUM('in_progress', 'close');

CREATE TABLE IF NOT EXISTS reception(
    ID uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    dateTime timestamp NOT NULL DEFAULT now(),
    pvzId uuid NOT NULL REFERENCES "PVZ"(id) ON DELETE CASCADE,
    status status NOT NULL DEFAULT 'in_progress'
)