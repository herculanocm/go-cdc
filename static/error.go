package static

type ErrorUtil interface {
	Error() string
	Code() string
	ToString() string
	SetPropagationError(err error)
	SetPropagationErrorString(errString string)
}

type AbstractError struct {
	error                  string
	code                   string
	propagationError       error
	propagationErrorString string
}

func (e *AbstractError) Error() string {
	return e.error
}

func (e *AbstractError) Code() string {
	return e.code
}

func (e *AbstractError) ToString() string {
	return "ErrorUtil { Code: " + e.code + ", Message: " + e.error + ", PropagationErrorString: " + e.propagationErrorString + " }"
}

func (e *AbstractError) SetPropagationError(err error) {
	e.propagationError = err
}

func (e *AbstractError) SetPropagationErrorString(errString string) {
	e.propagationErrorString = errString
}

func NewErrorUtil(errorMsg string, code string, propagationError error, propagationErrorString string) ErrorUtil {
	return &AbstractError{
		error:                  errorMsg,
		code:                   code,
		propagationError:       propagationError,
		propagationErrorString: propagationErrorString,
	}
}

var (
	ErrConfigLoadFailed        = NewErrorUtil("Failed to load configuration", "CONFIG_LOAD_FAILED", nil, "")
	ErrConfigLoadDecodeFailed  = NewErrorUtil("Failed to decode configuration", "CONFIG_DECODE_FAILED", nil, "")
	ErrConfigFileNotFound      = NewErrorUtil("Configuration file not found", "CONFIG_FILE_NOT_FOUND", nil, "")
	ErrEnvVarMissing           = NewErrorUtil("Required environment variable is missing", "ENV_VAR_MISSING", nil, "")
	ErrDBConnectionFailed      = NewErrorUtil("Failed to connect to the database", "DB_CONNECTION_FAILED", nil, "")
	ErrLoggerInitFailed        = NewErrorUtil("Failed to initialize logger", "LOGGER_INIT_FAILED", nil, "")
	ErrUnsupportedDBTechnology = NewErrorUtil("Unsupported database technology", "UNSUPPORTED_DB_TECHNOLOGY", nil, "")
)
