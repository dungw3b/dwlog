package dwlog

type Level uint32

const (
	UNKNOWN = iota
	ERROR
	INFO
	DEBUG
)

func LevelString(level uint32) string {
	switch level {
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

func (l Level) String() string {
	return LevelString(l.Val())
}

func (l Level) ColorString() string {
	switch l {
	case ERROR:
		return "\033[31mERROR\033[0m"
	case INFO:
		return "\033[32mINFO\033[0m"
	case DEBUG:
		return "\033[34mDEBUG\033[0m"
	default:
		return "UNKNOWN"
	}
}

func (l Level) Val() uint32 {
	return uint32(l)
}