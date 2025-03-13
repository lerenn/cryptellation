CREATE USER uptrace;
ALTER USER uptrace PASSWORD 'uptrace';
ALTER USER uptrace CREATEDB;

CREATE DATABASE uptrace;
GRANT ALL PRIVILEGES ON DATABASE uptrace TO uptrace;
\c uptrace postgres
GRANT ALL ON SCHEMA public TO uptrace;
