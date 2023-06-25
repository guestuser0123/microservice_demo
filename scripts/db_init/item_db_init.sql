CREATE DATABASE IF NOT EXISTS item_db;

USE item_db;

DROP TABLE IF EXISTS item_tab;

CREATE TABLE item_tab (
    item_id INT(11) AUTO_INCREMENT PRIMARY KEY,
    item_qty INT(11),
    item_name VARCHAR(20),
    CHECK(item_qty >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS item_audit_tab (
    txn_id VARCHAR(64) PRIMARY KEY,
    item_id INT(11),
    item_qty INT(11),
    status INT(11)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DELETE FROM item_audit_tab;

INSERT INTO item_tab (item_name, item_qty) VALUES ("ItemA", 10000), ("ItemB", 10000), ("ItemC", 10000), ("ItemC", 10000),
                                                  ("ItemD", 10000), ("ItemE", 10000), ("ItemF", 10000), ("ItemG", 10000),
                                                  ("ItemH", 10000), ("ItemI", 10000);
