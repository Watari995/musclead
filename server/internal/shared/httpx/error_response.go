package httpx

import (
	"errors"
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/getsentry/sentry-go"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string         `json:"code"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

func WriteError(w http.ResponseWriter, err error) {
	if myErr, ok := errors.AsType[myerror.MyError](err); ok {
		if myErr.Status() >= http.StatusInternalServerError {
			sentry.CaptureException(err)
		}
		WriteJSON(w, myErr.Status(), ErrorResponse{
			Error: ErrorDetail{
				Code:    string(myErr.Code()),
				Message: myErr.Message(),
				Data:    myErr.Data(),
			},
		})
		return
	}

	// それ以外ならinternal server error
	sentry.CaptureException(err)
	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Code: string(myerror.ErrorCodes.General.InternalError),
		},
	})

}
