package e2temperr

import (
	"errors"

	"github.com/m11ano/e"
	"go.temporal.io/sdk/temporal"
)

type TransportError struct {
	Code    int
	Message string
	Details []string
}

func ErrToTempErr(err error) error {
	if ok, lgErr := e.IsAppError(err); ok {
		nonRetryable := true

		if lgErr.Code() >= 500 && lgErr.Code() < 600 {
			nonRetryable = false
		}

		trErr := &TransportError{
			Code:    lgErr.Code(),
			Message: lgErr.Message(),
			Details: lgErr.Details(),
		}

		return temporal.NewApplicationErrorWithOptions(err.Error(), "WrapError", temporal.ApplicationErrorOptions{
			NonRetryable: nonRetryable,
			Details:      []any{trErr},
		})
	}

	return err
}

func TempErrConvertToLogicError(err error) (bool, *e.LogicError) {
	var appErr *temporal.ApplicationError
	if errors.As(err, &appErr) {
		if appErr.Type() == "" {
			return false, nil
		}

		var trErr *TransportError
		if appErr.HasDetails() {
			dErr := appErr.Details(&trErr)
			if dErr != nil {
				return false, nil
			}

			sourceErr := e.ErrInternal
			switch trErr.Code {
			case 400:
				sourceErr = e.ErrBadRequest
			case 401:
				sourceErr = e.ErrUnauthorized
			case 403:
				sourceErr = e.ErrForbidden
			case 404:
				sourceErr = e.ErrNotFound
			case 406:
				sourceErr = e.ErrNotAcceptable
			case 409:
				sourceErr = e.ErrConflict
			case 422:
				sourceErr = e.ErrUnprocessableEntity
			case 500:
				sourceErr = e.ErrInternal
			case 503:
				sourceErr = e.ErrServiceUnavailable
			}

			return true, e.NewErrorFrom(sourceErr).SetMessage(trErr.Message).AddDetails(trErr.Details)
		}
	}

	return false, nil
}
