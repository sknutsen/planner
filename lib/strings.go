package lib

import (
	"fmt"
	"strings"
)

// AsString coerces a sql/driver or JSON-decoded value to string. Nil becomes "".
func AsString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	default:
		return fmt.Sprint(s)
	}
}

func IsNullOrEmpty(s interface{}) bool {
	return s == nil || strings.TrimSpace(s.(string)) == ""
}
