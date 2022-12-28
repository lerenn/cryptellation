CREATE DATABASE IF NOT EXISTS ticks;

CREATE TABLE IF NOT EXISTS ticks.symbol_listeners
(
    exchange TEXT NOT NULL,
    pair_symbol TEXT NOT NULL,
    subscribers INT NOT NULL,
    PRIMARY KEY (exchange, pair_symbol)
);
