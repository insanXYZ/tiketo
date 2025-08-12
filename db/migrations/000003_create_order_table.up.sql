CREATE TYPE status AS ENUM ('unpaid','paid');

CREATE TABLE IF NOT EXISTS orders (
  id VARCHAR(100) NOT NULL,
  status status NOT NULL DEFAULT 'unpaid',
  total DECIMAL NOT NULL ,
  user_id VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  paid_at TIMESTAMP ,
  PRIMARY KEY(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
