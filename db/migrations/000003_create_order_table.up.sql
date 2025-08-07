CREATE TYPE status AS ENUM ('wait','paid','cancel');

CREATE TABLE IF NOT EXISTS orders (
  id VARCHAR(100) NOT NULL,
  status status NOT NULL DEFAULT 'wait',
  total MONEY NOT NULL ,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  paid_at TIMESTAMP ,
  PRIMARY KEY(id)
);
