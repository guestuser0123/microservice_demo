package repo

import (
	"database/sql"
	"fmt"
	"item_service/model"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	GetItemQuery    = `SELECT item_name, item_qty FROM item_tab WHERE item_id=?;`
	DeductItemQuery = `UPDATE item_tab SET item_qty=item_qty-? WHERE item_id=?;`
	AddItemQuery    = `UPDATE item_tab SET item_qty=item_qty+? WHERE item_id=?;`

	GetItemAuditQuery       = `SELECT item_id, item_qty FROM item_audit_tab WHERE txn_id=? AND status=?`
	AddItemAuditQuery       = `INSERT INTO item_audit_tab (txn_id, item_id, item_qty, status) VALUES (?, ?, ?, ?);`
	SetItemAuditStatusQuery = `UPDATE item_audit_tab SET status=? WHERE txn_id=?;`

	ItemAuditStatusBought   = 1
	ItemAuditStatusReturned = 2
)

type ItemRepo interface {
	MustInit()
	Close() error
	GetItemByID(itemID int64) (*model.Item, error)
	DeductItem(itemID, qty int64, txnID string) error
	ReturnItem(txnID string) error
}

type itemRepoImpl struct {
	db *sql.DB

	getItemStmt    *sql.Stmt
	deductItemStmt *sql.Stmt
	addItemStmt    *sql.Stmt

	getItemAuditStmt       *sql.Stmt
	addItemAuditStmt       *sql.Stmt
	setItemAuditStatusStmt *sql.Stmt
}

func NewRepo() ItemRepo {
	return &itemRepoImpl{}
}

func (r *itemRepoImpl) MustInit() {
	var (
		dbUser          = os.Getenv("ITEM_DB_USER")
		dbPwd           = os.Getenv("ITEM_DB_PWD")
		dbAddr          = os.Getenv("ITEM_DB_ADDR")
		maxIdleConns, _ = strconv.Atoi(os.Getenv("MYSQL_MAX_IDLE_CONNS"))
		maxOpenConns, _ = strconv.Atoi(os.Getenv("MYSQL_MAX_OPEN_CONNS"))
	)

	dbConnAddr := fmt.Sprintf("%v:%v@tcp(%v)/item_db", dbUser, dbPwd, dbAddr)
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
	r.getItemStmt = r.prepareQuery(GetItemQuery)
	r.addItemStmt = r.prepareQuery(AddItemQuery)
	r.deductItemStmt = r.prepareQuery(DeductItemQuery)
	r.getItemAuditStmt = r.prepareQuery(GetItemAuditQuery)
	r.addItemAuditStmt = r.prepareQuery(AddItemAuditQuery)
	r.setItemAuditStatusStmt = r.prepareQuery(SetItemAuditStatusQuery)
}

func (r *itemRepoImpl) prepareQuery(query string) *sql.Stmt {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	return stmt
}

func (r *itemRepoImpl) Close() error {
	return r.db.Close()
}

func (r *itemRepoImpl) GetItemByID(itemID int64) (*model.Item, error) {
	var (
		itemName string
		itemQty  int64
	)
	if err := r.getItemStmt.QueryRow(itemID).Scan(&itemName, &itemQty); err != nil {
		log.Println("GetItemById failed with error=", err.Error())
		return nil, err
	}
	return &model.Item{
		ID:   itemID,
		Qty:  itemQty,
		Name: itemName,
	}, nil
}

func (r *itemRepoImpl) DeductItem(itemID, qty int64, txnID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.Stmt(r.deductItemStmt).Exec(qty, itemID); err != nil {
		log.Println(err)
		return err
	}

	if _, err := tx.Stmt(r.addItemAuditStmt).Exec(txnID, itemID, qty, ItemAuditStatusBought); err != nil {
		log.Println(err)
		return err
	}

	return tx.Commit()
}

func (r *itemRepoImpl) ReturnItem(txnID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// 1. Check if audit exists and to get the qty to compensate. Else no action required.
	var itemID, qty int64
	if err = tx.Stmt(r.getItemAuditStmt).QueryRow(txnID, ItemAuditStatusBought).Scan(&itemID, &qty); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	// 2. Update item audit status.
	if _, err = tx.Stmt(r.setItemAuditStatusStmt).Exec(ItemAuditStatusReturned, txnID); err != nil {
		log.Println(err)
		return err
	}

	// 3. Return item quota.
	if _, err = tx.Stmt(r.addItemStmt).Exec(qty, itemID); err != nil {
		log.Println(err)
		return err
	}
	return tx.Commit()
}
