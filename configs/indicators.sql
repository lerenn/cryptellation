CREATE DATABASE IF NOT EXISTS indicators;

CREATE TABLE IF NOT EXISTS indicators.simple_moving_averages
(
    exchange_name TEXT NOT NULL,
    pair_symbol TEXT NOT NULL,
    period_symbol TEXT NOT NULL,
    period_number INT NOT NULL,
    price_type TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    price NUMERIC NOT NULL,
    PRIMARY KEY (exchange_name, pair_symbol, period_symbol, period_number, price_type, time)
);
