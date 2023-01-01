CREATE DATABASE IF NOT EXISTS backtests;

CREATE TABLE IF NOT EXISTS backtests.backtests
(
    id INT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    "current_time" TIMESTAMP NOT NULL,
    current_price_type TEXT NOT NULL,
    end_time TIMESTAMP NOT NULL,
    period_between_events TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS backtests.orders
(
    id INTEGER NOT NULL,
    execution_time TIMESTAMP,
    "type" TEXT NOT NULL, 
    backtest_id INT NOT NULL,
    exchange_name TEXT NOT NULL,
    pair_symbol TEXT NOT NULL,
    side TEXT NOT NULL,
    quantity NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_backtest_id FOREIGN KEY (backtest_id) REFERENCES backtests.backtests(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS backtests.tick_subscriptions
(
    id INTEGER NOT NULL,
    backtest_id INT NOT NULL,
    exchange_name TEXT NOT NULL,
    pair_symbol TEXT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_backtest_id FOREIGN KEY (backtest_id) REFERENCES backtests.backtests(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS backtests.balances
(
    "asset_name" TEXT NOT NULL,
    exchange_name TEXT NOT NULL,
    backtest_id INT NOT NULL,
    balance NUMERIC NOT NULL,
    PRIMARY KEY (backtest_id, exchange_name, "asset_name"),
    CONSTRAINT fk_backtest_id FOREIGN KEY (backtest_id) REFERENCES backtests.backtests(id) ON DELETE CASCADE ON UPDATE CASCADE
);
