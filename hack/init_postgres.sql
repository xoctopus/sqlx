CREATE DATABASE test;

\c test;

CREATE TABLE x (datetime timestamptz(3), timestamp bigint);

INSERT INTO x(datetime, timestamp) VALUES (CURRENT_TIMESTAMP, FLOOR(EXTRACT(EPOCH FROM now()) * 1000));