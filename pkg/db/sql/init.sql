BEGIN;
  CREATE TABLE IF NOT EXISTS sales(
    id int PRIMARY KEY,
    type int,
    date timestamp without time zone,
    product text,
    value real,
    saler text
  );
COMMIT;
