package myerror

import (
	"errors"
	"net/http"
)

type MyError interface {
	Error() string
	Unwrap() error
	Wrap(error) MyError
	Message() string
	SetMessage(message string) MyError
	Code() ErrorCode
	SetCode(code ErrorCode) MyError
	Status() int
	Data() map[string]any
	SetData(data ...map[string]any) MyError
}

type myError struct {
	generalCode ErrorCode
	message     string
	code        ErrorCode
	status      int
	data        map[string]any
	cause       error
}

func (e *myError) Error() string {
	if e.message != "" {
		return e.message
	}
	if e.code == "" {
		return string(e.generalCode)
	}
	return string(e.code)
}

func (e *myError) Unwrap() error {
	return e.cause
}

func (e *myError) Wrap(err error) MyError {
	e.cause = err
	return e
}

func (e *myError) Message() string {
	return e.message
}

func (e *myError) SetMessage(message string) MyError {
	e.message = message
	return e
}

func (e *myError) Code() ErrorCode {
	if e.code != "" {
		return e.code
	}
	return e.generalCode
}

func (e *myError) SetCode(code ErrorCode) MyError {
	e.code = code
	return e
}

func (e *myError) Status() int {
	return e.status
}

func (e *myError) Data() map[string]any {
	return e.data
}

func (e *myError) SetData(data ...map[string]any) MyError {
	if e.data == nil {
		e.data = make(map[string]any)
	}
	for _, d := range data {
		for k, v := range d {
			e.data[k] = v
		}
	}
	return e
}

// ─── 汎用ファクトリ ───

func NewInternalError() MyError {
	return &myError{
		status:      http.StatusInternalServerError,
		generalCode: ErrorCodes.General.InternalError,
	}
}

func NewValidationError() MyError {
	return &myError{
		status:      http.StatusBadRequest,
		generalCode: ErrorCodes.General.ValidationError,
	}
}

func NewPermissionError() MyError {
	return &myError{
		status:      http.StatusForbidden,
		generalCode: ErrorCodes.General.PermissionError,
	}
}

func NewBadRequestError() MyError {
	return &myError{
		status:      http.StatusBadRequest,
		generalCode: ErrorCodes.General.BadRequestError,
	}
}

func NewIncompatibleError() MyError {
	return &myError{
		status:      http.StatusUpgradeRequired,
		generalCode: ErrorCodes.General.IncompatibleError,
	}
}

func NewNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
	}
}

func NewUnauthorizedError() MyError {
	return &myError{
		status:      http.StatusUnauthorized,
		generalCode: ErrorCodes.General.UnauthorizedError,
	}
}

func NewMaintenanceModeError() MyError {
	return &myError{
		status:      http.StatusServiceUnavailable,
		generalCode: ErrorCodes.General.MaintenanceModeError,
	}
}

func NewDisabledUserError() MyError {
	return &myError{
		status:      http.StatusForbidden,
		generalCode: ErrorCodes.General.DisabledUserError,
	}
}

// ─── ドメイン固有ファクトリ (generalCode + code を分離) ───

func NewUserNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.User.NotFoundError,
	}
}

func NewEmailAlreadyExistsError() MyError {
	return &myError{
		status:      http.StatusConflict,
		generalCode: ErrorCodes.General.BadRequestError,
		code:        ErrorCodes.User.EmailAlreadyExistsError,
	}
}

func NewInvalidCredentialsError() MyError {
	return &myError{
		status:      http.StatusUnauthorized,
		generalCode: ErrorCodes.General.UnauthorizedError,
		code:        ErrorCodes.User.InvalidCredentialsError,
	}
}

func NewMealNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.Meal.NotFoundError,
	}
}

func NewTrainingNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.Training.NotFoundError,
	}
}

func NewExerciseNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.Training.ExerciseNotFoundError,
	}
}

func NewExerciseNameAlreadyExistsError() MyError {
	return &myError{
		status:      http.StatusConflict,
		generalCode: ErrorCodes.General.BadRequestError,
		code:        ErrorCodes.Training.ExerciseNameAlreadyExistsError,
	}
}

func NewExerciseUsedInTrainingError() MyError {
	return &myError{
		status:      http.StatusConflict,
		generalCode: ErrorCodes.General.BadRequestError,
		code:        ErrorCodes.Training.ExerciseUsedInTrainingError,
	}
}

func NewRoutineNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.Training.RoutineNotFoundError,
	}
}

func NewRoutineNameAlreadyExistsError() MyError {
	return &myError{
		status:      http.StatusConflict,
		generalCode: ErrorCodes.General.BadRequestError,
		code:        ErrorCodes.Training.RoutineNameAlreadyExistsError,
	}
}

func NewWeightNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.Weight.NotFoundError,
	}
}

func NewSubscriptionOrderNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.SubscriptionOrder.NotFoundError,
	}
}

func NewPaymentNotFoundError() MyError {
	return &myError{
		status:      http.StatusNotFound,
		generalCode: ErrorCodes.General.NotFoundError,
		code:        ErrorCodes.Payment.NotFoundError,
	}
}

// ─── ヘルパー ───
func IsCode(err error, code ErrorCode) bool {
	if myErr, ok := errors.AsType[MyError](err); ok {
		return myErr.Code() == code
	}
	return false
}
