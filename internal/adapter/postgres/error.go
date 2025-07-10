package postgres

import "errors"

var (
	ErrNotFound    = errors.New("data not found, try again later")
	ErrQueryFailed = errors.New("query failed")
	ErrScanFailed  = errors.New("scan failed")
)
