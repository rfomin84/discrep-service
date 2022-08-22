CREATE TABLE balance_history (
    id INT NOT NULL AUTO_INCREMENT,
    feed_id INT NOT NULL,
    date DATETIME NOT NULL,
    cost INT,
    approved TINYINT DEFAULT 0,
    PRIMARY KEY (id),
    INDEX (feed_id)
);