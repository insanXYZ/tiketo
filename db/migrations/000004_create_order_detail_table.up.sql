CREATE TABLE IF NOT EXISTS order_details (
  id SERIAL,
  order_id VARCHAR(100) NOT NULL,
  ticket_id VARCHAR(100) NOT NULL,
  quantity INT NOT NULL, 
  PRIMARY KEY(id)
);
