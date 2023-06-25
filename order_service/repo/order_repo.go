package repo

import (
	"database/sql"
	"fmt"
	"log"
	"order_service/model"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	AddOrderQuery = `INSERT INTO order_tab (txn_id, item_id, item_qty, voucher_id) VALUES (?, ?, ?, ?);`
)

type OrderRepo interface {
	MustInit()
	Close() error
	AddOrder(txnID string, itemID, itemQty, voucherID int64) (*model.Order, error)
}

type orderRepoImpl struct {
	db *sql.DB

	AddOrderStmt *sql.Stmt
}

func NewRepo() OrderRepo {
	return &orderRepoImpl{}
}

func (r *orderRepoImpl) MustInit() {
	var (
		dbUser          = os.Getenv("ORDER_DB_USER")
		dbPwd           = os.Getenv("ORDER_DB_PWD")
		dbAddr          = os.Getenv("ORDER_DB_ADDR")
		maxIdleConns, _ = strconv.Atoi(os.Getenv("MYSQL_MAX_IDLE_CONNS"))
		maxOpenConns, _ = strconv.Atoi(os.Getenv("MYSQL_MAX_OPEN_CONNS"))
	)

	dbConnAddr := fmt.Sprintf("%v:%v@tcp(%v)/order_db", dbUser, dbPwd, dbAddr)
	dbConn, err := sql.Open("mysql", dbConnAddr)
	if err != nil {
		log.Fatal(err)
	}

	dbConn.SetMaxIdleConns(maxIdleConns)
	dbConn.SetMaxOpenConns(maxOpenConns)
	for i := 0; i < maxIdleConns; i++ {
		dbConn.Ping()
	}

	r.db = dbConn
	r.AddOrderStmt = r.prepareQuery(AddOrderQuery)
}

func (r *orderRepoImpl) prepareQuery(query string) *sql.Stmt {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	return stmt
}

func (r *orderRepoImpl) Close() error {
	return r.db.Close()
}

func (r *orderRepoImpl) AddOrder(txnID string, itemID, itemQty, voucherID int64) (*model.Order, error) {
	_, err := r.AddOrderStmt.Exec(txnID, itemID, itemQty, voucherID)
	if err != nil {
		return nil, err
	}
	return &model.Order{
		TxnID:     txnID,
		ItemID:    itemID,
		ItemQty:   itemQty,
		VoucherID: voucherID,
	}, nil
}
