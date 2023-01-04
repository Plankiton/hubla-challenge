BEGIN;
  CREATE TABLE sales(
    id SERIAL PRIMARY KEY,
    type int,
    date timestamp without time zone,
    product text,
    value real,
    saler text
  );
COMMIT;
