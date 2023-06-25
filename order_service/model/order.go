package model

type Order struct {
	TxnID     string
	ItemID    int64
	ItemQty   int64
	VoucherID int64
}
