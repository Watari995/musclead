package myerror

type ErrorCode string

type generalErrors struct {
	InternalError        ErrorCode
	ValidationError      ErrorCode
	PermissionError      ErrorCode
	BadRequestError      ErrorCode
	IncompatibleError    ErrorCode
	NotFoundError        ErrorCode
	UnauthorizedError    ErrorCode
	MaintenanceModeError ErrorCode
	DisabledUserError    ErrorCode
}

type userErrors struct {
	NotFoundError           ErrorCode
	EmailAlreadyExistsError ErrorCode
	InvalidCredentialsError ErrorCode
}

type mealErrors struct {
	NotFoundError ErrorCode
}

type trainingErrors struct {
	NotFoundError                  ErrorCode
	ExerciseNotFoundError          ErrorCode
	ExerciseNameAlreadyExistsError ErrorCode
}

type weightErrors struct {
	NotFoundError ErrorCode
}

var ErrorCodes = struct {
	General  generalErrors
	User     userErrors
	Meal     mealErrors
	Training trainingErrors
	Weight   weightErrors
}{
	General: generalErrors{
		InternalError:        "general.internal_error",
		ValidationError:      "general.validation_error",
		PermissionError:      "general.permission_error",
		BadRequestError:      "general.bad_request_error",
		IncompatibleError:    "general.incompatible_error",
		NotFoundError:        "general.not_found_error",
		UnauthorizedError:    "general.unauthorized_error",
		MaintenanceModeError: "general.maintenance_mode_error",
		DisabledUserError:    "general.disabled_user_error",
	},
	User: userErrors{
		NotFoundError:           "user.not_found_error",
		EmailAlreadyExistsError: "user.email_already_exists_error",
		InvalidCredentialsError: "user.invalid_credentials_error",
	},
	Meal: mealErrors{
		NotFoundError: "meal.not_found_error",
	},
	Training: trainingErrors{
		NotFoundError:                  "training.not_found_error",
		ExerciseNotFoundError:          "training.exercise_not_found_error",
		ExerciseNameAlreadyExistsError: "training.exercise_name_already_exists_error",
	},
	Weight: weightErrors{
		NotFoundError: "weight.not_found_error",
	},
}
