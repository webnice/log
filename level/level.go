package level // import "github.com/webnice/log/v2/level"

// New Create new level object
func New() Interface {
	var level = new(impl)
	return level
}

// -------
// Interface
// -------

// Fatal system is unusable
func (l *impl) Fatal() Level { return levelFatal }

// Alert action must be taken immediately
func (l *impl) Alert() Level { return levelAlert }

// Critical conditions
func (l *impl) Critical() Level { return levelCritical }

// Error conditions
func (l *impl) Error() Level { return levelError }

// Warning conditions
func (l *impl) Warning() Level { return levelWarning }

// Notice normal but significant condition
func (l *impl) Notice() Level { return levelNotice }

// Informational messages
func (l *impl) Informational() Level { return levelInfo }

// Debug debug-level messages
func (l *impl) Debug() Level { return levelDebug }

// -------
// Level
// -------

// String Return Level as string
func (lvl Level) String() string {
	return levelMap[lvl]
}

// Int8 Return Level as int8
func (lvl Level) Int8() int8 {
	return int8(lvl)
}

// Int Return Level as int
func (lvl Level) Int() int {
	return int(lvl)
}

// Int64 Return Level as int64
func (lvl Level) Int64() int64 {
	return int64(lvl)
}
