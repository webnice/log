package level

// New Create new level object
func New() Interface {
	var level = new(impl)
	return level
}

// -------
// Interface
// -------

// Fatal system is unusable
func (l *impl) Fatal() Level { return _FATAL }

// Alert action must be taken immediately
func (l *impl) Alert() Level { return _ALERT }

// Critical conditions
func (l *impl) Critical() Level { return _CRITICAL }

// Error conditions
func (l *impl) Error() Level { return _ERROR }

// Warning conditions
func (l *impl) Warning() Level { return _WARNING }

// Notice normal but significant condition
func (l *impl) Notice() Level { return _NOTICE }

// Informational messages
func (l *impl) Informational() Level { return _INFO }

// Debug debug-level messages
func (l *impl) Debug() Level { return _DEBUG }

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
