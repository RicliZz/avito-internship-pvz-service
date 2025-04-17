CREATE TABLE IF NOT EXISTS "PVZ" (
    "ID" uuid PRIMARY KEY default gen_random_uuid(),
    "registrationDate" TIMESTAMP NOT NULL DEFAULT now(),
    city TEXT NOT NULL CHECK ( city IN ('Москва', 'Санкт-Петербург', 'Казань') )
)