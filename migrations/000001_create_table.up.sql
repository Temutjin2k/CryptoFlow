CREATE TABLE IF NOT EXISTS aggregated_prices (
    id SERIAL PRIMARY KEY,
    pair_name TEXT NOT NULL,
    exchange TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    min_price FLOAT NOT NULL,
    max_price FLOAT NOT NULL,
    average_price FLOAT NOT NULL
);