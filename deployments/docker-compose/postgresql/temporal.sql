CREATE USER temporal;
ALTER USER temporal PASSWORD 'temporal';
ALTER USER temporal CREATEDB;

CREATE DATABASE temporal;
GRANT ALL PRIVILEGES ON DATABASE temporal TO temporal;
\c temporal postgres
GRANT ALL ON SCHEMA public TO temporal;
