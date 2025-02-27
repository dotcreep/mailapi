package utils

import (
	"errors"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func GenerateBarCodeImage(url string) (barcode.Barcode, error) {
	// TODO: Implement barcode generation logic here
	if url == "" {
		return nil, errors.New("url is required")
	}
	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}
