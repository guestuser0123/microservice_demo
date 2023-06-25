package repo

import (
	"errors"
	"log"
	"net/http"
	"os"
)

type ItemRepo interface {
	DeductItem(itemID, itemQty, txnID string) error
	DeductItemCancel(txnID string) error
}

type itemRepo struct {
	clientAdapter *http.Client
}

var (
	ItemServiceUrl = "http://" + os.Getenv("ITEM_SERVICE_INTERNAL_ADDR")

	DeductItemAPIPath = "/item/deduct"
	CancelItemAPIPath = "/item/cancel"
)

func NewItemRepo(clientAdapter *http.Client) ItemRepo {
	return &itemRepo{
		clientAdapter: clientAdapter,
	}
}

func (r *itemRepo) DeductItem(itemID, itemQty, txnID string) error {
	url := ItemServiceUrl + DeductItemAPIPath + "?itemid=" + itemID + "&itemqty=" + itemQty + "&txnid=" + txnID
	req, err := http.NewRequest(http.MethodPost, url, nil)
	res, err := r.clientAdapter.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("failed to deduct item")
	}

	return nil
}

func (r *itemRepo) DeductItemCancel(txnID string) error {
	url := ItemServiceUrl + CancelItemAPIPath + "?txnid=" + txnID
	req, err := http.NewRequest(http.MethodPost, url, nil)
	res, err := r.clientAdapter.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("failed to cancel item")
	}

	return nil
}
