package repo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"voucher_service/model"

	_ "github.com/go-sql-driver/mysql"
)

const (
	GetVoucherQuery            = `SELECT voucher_name, voucher_qty FROM voucher_tab WHERE voucher_id=?;`
	DeductVoucherQtyQuery      = `UPDATE voucher_tab SET voucher_qty=voucher_qty-? WHERE voucher_id=?;`
	AddVoucherQtyQuery         = `UPDATE voucher_tab SET voucher_qty=voucher_qty+? WHERE voucher_id=?;`
	GetVoucherAuditQuery       = `SELECT voucher_id, voucher_qty FROM voucher_audit_tab WHERE txn_id=? and status=?`
	AddVoucherAuditQuery       = `INSERT INTO voucher_audit_tab (txn_id, voucher_id, voucher_qty, status) VALUES (?, ?, ?, ?);`
	SetVoucherAuditStatusQuery = `UPDATE voucher_audit_tab SET status=? WHERE txn_id=?;`

	VoucherAuditStatusUsed     = 1
	VoucherAuditStatusReturned = 2
)

type VoucherRepo interface {
	MustInit()
	Close() error
	GetVoucherByID(voucherID int64) (*model.Voucher, error)
	UseVoucher(voucherID, qty int64, txnID string) error
	ReturnVoucher(txnID string) error
}

type voucherRepoImpl struct {
	db      *sql.DB
	closeFn func() error

	getVoucherStmt            *sql.Stmt
	getVoucherAuditStmt       *sql.Stmt
	deductVoucherQtyStmt      *sql.Stmt
	addVoucherQtyStmt         *sql.Stmt
	addVoucherAuditStmt       *sql.Stmt
	setVoucherAuditStatusStmt *sql.Stmt
}

func NewVoucherRepo() VoucherRepo {
	return &voucherRepoImpl{}
}

func (r *voucherRepoImpl) MustInit() {
	var (
		dbUser          = os.Getenv("VOUCHER_DB_USER")
		dbPwd           = os.Getenv("VOUCHER_DB_PWD")
		dbAddr          = os.Getenv("VOUCHER_DB_ADDR")
		maxIdleConns, _ = strconv.Atoi(os.Getenv("MYSQL_MAX_IDLE_CONNS"))
		maxOpenConns, _ = strconv.Atoi(os.Getenv("MYSQL_MAX_OPEN_CONNS"))
	)

	dbConnAddr := fmt.Sprintf("%v:%v@tcp(%v)/voucher_db", dbUser, dbPwd, dbAddr)
	dbConn, err := sql.Open("mysql", dbConnAddr)
	if err != nil {
		log.Fatalln(err)
	}

	dbConn.SetMaxIdleConns(maxIdleConns)
	dbConn.SetMaxOpenConns(maxOpenConns)
	for i := 0; i < maxIdleConns; i++ {
		dbConn.Ping()
	}

	r.db = dbConn
	r.addVoucherAuditStmt = r.prepareQuery(AddVoucherAuditQuery)
	r.deductVoucherQtyStmt = r.prepareQuery(DeductVoucherQtyQuery)
	r.addVoucherQtyStmt = r.prepareQuery(AddVoucherQtyQuery)
	r.setVoucherAuditStatusStmt = r.prepareQuery(SetVoucherAuditStatusQuery)
	r.getVoucherStmt = r.prepareQuery(GetVoucherQuery)
	r.getVoucherAuditStmt = r.prepareQuery(GetVoucherAuditQuery)
}

func (r *voucherRepoImpl) Close() error {
	return r.db.Close()
}

func (r *voucherRepoImpl) prepareQuery(query string) *sql.Stmt {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Fatalln(err)
	}
	return stmt
}

func (r *voucherRepoImpl) GetVoucherByID(voucherID int64) (*model.Voucher, error) {
	var (
		voucherName string
		voucherQty  int64
	)
	if err := r.getVoucherStmt.QueryRow(voucherID).Scan(&voucherName, &voucherQty); err != nil {
		log.Println("GetVoucherById failed with error=", err.Error())
		return nil, err
	}
	return &model.Voucher{
		ID:   voucherID,
		Qty:  voucherQty,
		Name: voucherName,
	}, nil
}

func (r *voucherRepoImpl) UseVoucher(voucherID, qty int64, txnID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.Stmt(r.deductVoucherQtyStmt).Exec(qty, voucherID); err != nil {
		log.Println(err)
		return err
	}
	if _, err = tx.Stmt(r.addVoucherAuditStmt).Exec(txnID, voucherID, qty, VoucherAuditStatusUsed); err != nil {
		log.Println(err)
		return err
	}
	return tx.Commit()
}

func (r *voucherRepoImpl) ReturnVoucher(txnID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// 1. Check if audit exists and to get the qty to compensate. Else no action required.
	var voucherID, qty int64
	if err = tx.Stmt(r.getVoucherAuditStmt).QueryRow(txnID, VoucherAuditStatusUsed).Scan(&voucherID, &qty); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	// 2. Update voucher audit status.
	if _, err = tx.Stmt(r.setVoucherAuditStatusStmt).Exec(VoucherAuditStatusReturned, txnID); err != nil {
		log.Println(err)
		return err
	}

	// 3. Return voucher quota.
	if _, err = tx.Stmt(r.addVoucherQtyStmt).Exec(qty, voucherID); err != nil {
		log.Println(err)
		return err
	}
	return tx.Commit()
}
