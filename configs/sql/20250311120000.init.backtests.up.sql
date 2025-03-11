CREATE TABLE backtests
(
    id VARCHAR(255) NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_backtests PRIMARY KEY (id)
)