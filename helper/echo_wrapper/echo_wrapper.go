package echo_wrapper

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"

	pkgError "github.com/adityaeka26/go-pkg/error"
)

// Result common output
type Result struct {
	Data     any
	MetaData any
	Error    error
	Count    int64
}

type response struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func getErrorStatusCode(err error) int {
	errString, ok := err.(*pkgError.ErrorString)
	if ok {
		return errString.Code()
	}

	return http.StatusInternalServerError
}

func RespSuccess(c echo.Context, data any, message string) error {
	return c.JSON(http.StatusOK, response{
		Message: message,
		Data:    data,
		Code:    http.StatusOK,
		Success: true,
	})
}

func RespError(c echo.Context, err error) error {
	return c.JSON(getErrorStatusCode(err), response{
		Message: err.Error(),
		Data:    nil,
		Code:    getErrorStatusCode(err),
		Success: false,
	})
}

type MetaData struct {
	Page      int64 `json:"page"`
	Count     int64 `json:"count"`
	TotalPage int64 `json:"total_page"`
	TotalData int64 `json:"total_data"`
}

func GenerateMetaData(totalData, count int64, page, limit int64) MetaData {
	metaData := MetaData{
		Page:      page,
		Count:     count,
		TotalPage: int64(math.Ceil(float64(totalData) / float64(limit))),
		TotalData: totalData,
	}

	return metaData
}

type paginationResponse struct {
	Success bool     `json:"success"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Meta    MetaData `json:"meta"`
	Data    any      `json:"data"`
}

func RespPagination(c echo.Context, data any, metadata MetaData, message string) error {
	return c.JSON(http.StatusOK, paginationResponse{
		Message: message,
		Meta:    metadata,
		Data:    data,
		Code:    http.StatusOK,
		Success: true,
	})
}
