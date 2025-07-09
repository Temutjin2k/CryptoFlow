CREATE TABLE IF NOT EXISTS latest_prices(
    exchange VARCHAR(100) NOT NULL,
    pair_name VARCHAR NOT NULL,
    price FLOAT NOT NULL,
    timestamp Timestamp NOT NULL,
    PRIMARY KEY (exchange, pair_name)
)