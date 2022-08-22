CREATE TABLE discrepancy_statistics (
  id INT NOT NULL AUTO_INCREMENT,
  feed_id INT NOT NULL,
  date DATE NOT NULL,
  cost INT,
  external_cost INT,
  discrepancy DECIMAL(5, 4),
  finalized TINYINT DEFAULT 0,
  PRIMARY KEY (id)
);
