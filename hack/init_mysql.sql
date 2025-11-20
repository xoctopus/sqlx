CREATE DATABASE test;

USE test;

CREATE TABLE x (datetime datetime(3), timestamp bigint);

INSERT INTO x(datetime, timestamp) VALUES (now(3), unix_timestamp(now(3)) * 1000);
