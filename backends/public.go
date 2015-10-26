package backends

import (
	"strings"
)

func CheckMode(m BackendName) BackendName {
	var mode BackendName
	mode = BackendName(strings.ToLower(string(m)))
	return mode
}
