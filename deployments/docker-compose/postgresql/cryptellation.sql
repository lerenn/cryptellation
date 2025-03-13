CREATE USER cryptellation;
ALTER USER cryptellation PASSWORD 'cryptellation';
ALTER USER cryptellation CREATEDB;

CREATE DATABASE cryptellation;
GRANT ALL PRIVILEGES ON DATABASE cryptellation TO cryptellation;
\c cryptellation postgres
GRANT ALL ON SCHEMA public TO cryptellation;