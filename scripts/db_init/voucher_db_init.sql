CREATE DATABASE IF NOT EXISTS voucher_db;

USE voucher_db;

DROP TABLE IF EXISTS voucher_tab;

CREATE TABLE voucher_tab (
    voucher_id INT(11) AUTO_INCREMENT PRIMARY KEY,
    voucher_qty INT(11),
    voucher_name VARCHAR(20),
    CHECK(voucher_qty >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS voucher_audit_tab (
    txn_id VARCHAR(64) PRIMARY KEY,
    voucher_id INT(11),
    voucher_qty INT(11),
    status INT(11)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DELETE FROM voucher_audit_tab;

INSERT INTO voucher_tab (voucher_name, voucher_qty) VALUES ("VoucherA", 10000), ("VoucherB", 10000), ("VoucherC", 10000), ("VoucherC", 10000),
                                                           ("VoucherD", 10000), ("VoucherE", 10000), ("VoucherF", 10000), ("VoucherG", 10000),
                                                           ("VoucherH", 10000), ("VoucherI", 10000);
