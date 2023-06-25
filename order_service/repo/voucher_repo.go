package repo

import (
	"errors"
	"log"
	"net/http"
	"os"
)

type VoucherRepo interface {
	UseVoucher(voucherID, txnID string) error
	UseVoucherCancel(txnID string) error
}

type voucherRepo struct {
	clientAdapter *http.Client
}

func NewVoucherRepo(clientAdapter *http.Client) VoucherRepo {
	return &voucherRepo{
		clientAdapter: clientAdapter,
	}
}

var (
	VoucherServiceUrl = "http://" + os.Getenv("VOUCHER_SERVICE_INTERNAL_ADDR")

	UseVoucherAPIPath    = "/voucher/use"
	CancelVoucherAPIPath = "/voucher/cancel"
)

func (r *voucherRepo) UseVoucher(voucherID, txnID string) error {
	url := VoucherServiceUrl + UseVoucherAPIPath + "?voucherid=" + voucherID + "&txnid=" + txnID
	req, err := http.NewRequest(http.MethodPost, url, nil)
	res, err := r.clientAdapter.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("failed to use voucher")
	}

	return nil
}

func (r *voucherRepo) UseVoucherCancel(txnID string) error {
	url := VoucherServiceUrl + CancelVoucherAPIPath + "?txnid=" + txnID
	req, err := http.NewRequest(http.MethodPost, url, nil)
	res, err := r.clientAdapter.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("failed to cancel voucher")
	}

	return nil
}
