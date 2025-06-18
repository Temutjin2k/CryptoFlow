CREATE TABLE IF NOT EXISTS aggregated_prices (
    id SERIAL PRIMARY KEY,
    pair_name TEXT NOT NULL,
    exchange TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,  -- minute-level rounded timestamp
    min_price FLOAT8 NOT NULL,
    max_price FLOAT8 NOT NULL,
    average_price FLOAT8 NOT NULL
);