package dwlog

type Level uint

const (
	ERROR = iota
	INFO
	DEBUG
)

func (l Level) String() string {
	switch l {
	case ERROR:
		return "ERROR"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}