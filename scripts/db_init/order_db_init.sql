CREATE DATABASE IF NOT EXISTS order_db;

USE order_db;

CREATE TABLE IF NOT EXISTS order_tab (
    txn_id VARCHAR(64) PRIMARY KEY,
    item_id INT(11),
    item_qty INT(11),
    voucher_id INT(11)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DELETE FROM order_tab;
