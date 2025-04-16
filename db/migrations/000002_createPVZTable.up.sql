CREATE TABLE IF NOT EXISTS "PVZ" (
    id uuid PRIMARY KEY default gen_random_uuid(),
    registration_date TIMESTAMP NOT NULL DEFAULT now(),
    city TEXT NOT NULL CHECK ( city IN ('Москва', 'Санкт-Петербург', 'Казань') )
)