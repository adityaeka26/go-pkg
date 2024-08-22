package error

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertHttpCodeToGrpcCode(err error) codes.Code {
	errString, ok := err.(*ErrorString)
	httpCode := http.StatusInternalServerError
	if ok {
		httpCode = errString.Code()
	}

	httpToGrpcCode := map[int]codes.Code{
		http.StatusOK:                  codes.OK,
		499:                            codes.Canceled,
		http.StatusInternalServerError: codes.Unknown,
		http.StatusBadRequest:          codes.InvalidArgument,
		http.StatusGatewayTimeout:      codes.DeadlineExceeded,
		http.StatusNotFound:            codes.NotFound,
		http.StatusConflict:            codes.AlreadyExists,
		http.StatusForbidden:           codes.PermissionDenied,
		http.StatusTooManyRequests:     codes.ResourceExhausted,
		http.StatusNotImplemented:      codes.Unimplemented,
		http.StatusServiceUnavailable:  codes.Unavailable,
	}
	if grpcCode, exists := httpToGrpcCode[httpCode]; exists {
		return grpcCode
	}

	return codes.Internal
}

func convertGrpcCodeToHttpCode(grpcCode codes.Code) int {
	grpcToHttpCode := map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           499, // Client Closed Request
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusGatewayTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
	}

	if httpCode, exists := grpcToHttpCode[grpcCode]; exists {
		return httpCode
	}

	return http.StatusInternalServerError
}

func GrpcError(err error) error {
	return status.Error(convertHttpCodeToGrpcCode(err), err.Error())
}

func HttpError(err error) error {
	return &ErrorString{
		code:     convertGrpcCodeToHttpCode(status.Code(err)),
		message:  status.Convert(err).Message(),
		httpCode: convertGrpcCodeToHttpCode(status.Code(err)),
	}
}
