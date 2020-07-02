package constants

import (
	"crawunit/errors"
	"net/http"
)

var (
	ErrorMongoError = errors.NewSystemOPError(1001, "system error:[mongo]%s", http.StatusInternalServerError)
	ErrorDepError = errors.NewSystemOPError(1002, "system error:[dep]%s", http.StatusInternalServerError)
	ErrorTaskExecuteError = errors.NewSystemOPError(1003, "system error:[task execute]%s", http.StatusInternalServerError)

)
