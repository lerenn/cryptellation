CREATE TABLE exchanges
(
    name VARCHAR(100) NOT NULL,
    data JSONB NOT NULL,
    CONSTRAINT pk_exchanges PRIMARY KEY (name)
)