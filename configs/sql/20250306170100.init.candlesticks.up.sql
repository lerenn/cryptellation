CREATE TABLE candlesticks
(
    exchange VARCHAR(100) NOT NULL,
    pair VARCHAR(100) NOT NULL,
    period VARCHAR(100) NOT NULL,
    time TIMESTAMP NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_candlesticks PRIMARY KEY (exchange, pair, period, time)
)