CREATE TABLE forwardtests
(
    id VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP,
    data JSONB NOT NULL,
    CONSTRAINT pk_forwardtests PRIMARY KEY (id)
)