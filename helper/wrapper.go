package helper

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	pkgError "adityaeka26/go-pkg/error"
)

func init() {}

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

type Meta struct {
	Method        string    `json:"method"`
	Url           string    `json:"url"`
	Code          string    `json:"code"`
	ContentLength int64     `json:"content_length"`
	Date          time.Time `json:"date"`
	Ip            string    `json:"ip"`
}

func RespSuccess(c *fiber.Ctx, data any, message string) error {
	return c.Status(http.StatusOK).JSON(response{
		Message: message,
		Data:    data,
		Code:    http.StatusOK,
		Success: true,
	})
}

func RespError(c *fiber.Ctx, err error) error {
	return c.Status(getErrorStatusCode(err)).JSON(response{
		Message: err.Error(),
		Data:    nil,
		Code:    getErrorStatusCode(err),
		Success: false,
	})
}
