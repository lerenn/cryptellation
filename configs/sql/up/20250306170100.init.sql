CREATE TABLE candlesticks
(
    exchange VARCHAR(100) NOT NULL,
    pair VARCHAR(100) NOT NULL,
    period VARCHAR(100) NOT NULL,
    time TIMESTAMP NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_candlesticks PRIMARY KEY (exchange, pair, period, time)
);

CREATE TABLE exchanges
(
    name VARCHAR(100) NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_exchanges PRIMARY KEY (name)
);

CREATE TABLE backtests
(
    id VARCHAR(255) NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_backtests PRIMARY KEY (id)
);

CREATE TABLE forwardtests
(
    id VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP,
    data JSONB NOT NULL,
    CONSTRAINT pk_forwardtests PRIMARY KEY (id)
);

CREATE TABLE indicators_sma
(
    exchange VARCHAR(100) NOT NULL,
    pair VARCHAR(100) NOT NULL,
    period VARCHAR(100) NOT NULL,
    period_number INTEGER NOT NULL,
    price_type VARCHAR(100) NOT NULL,
    time TIMESTAMP NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_indicators_sma PRIMARY KEY (exchange, pair, period, period_number, price_type, time)
);