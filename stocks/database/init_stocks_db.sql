DROP DATABASE IF EXISTS stocks;
DROP ROLE IF EXISTS user_stocks;

CREATE ROLE user_stocks LOGIN PASSWORD 'stocks1234';

CREATE DATABASE stocks;
    WITH OWNER = user_stocks
				 ENCODING = 'UTF8'
				 CONNECTION LIMIT = -1
				 IS_TEMPLATE = False;


\c stocks


GRANT ALL PRIVILEGES ON DATABASE stocks TO user_stocks;
GRANT USAGE ON SCHEMA public TO user_stocks;
GRANT CREATE ON SCHEMA public TO user_stocks;
