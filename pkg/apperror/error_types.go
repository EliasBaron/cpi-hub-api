package apperror

// ErrorType represents the type of error
type ErrorType int

const (
	NotFound ErrorType = iota
	UnexpectedDatabaseError
	Forbidden
	TooManyRequests
	Gone
	InternalServer
	InvalidData
	PayloadTooLarge
)

// String returns the string representation of the error type
func (e ErrorType) String() string {
	switch e {
	case NotFound:
		return "NotFound"
	case UnexpectedDatabaseError:
		return "UnexpectedDatabaseError"
	case Forbidden:
		return "Forbidden"
	case TooManyRequests:
		return "TooManyRequests"
	case Gone:
		return "Gone"
	case InternalServer:
		return "InternalServer"
	case InvalidData:
		return "InvalidData"
	case PayloadTooLarge:
		return "PayloadTooLarge"
	default:
		return "Unknown"
	}
}
