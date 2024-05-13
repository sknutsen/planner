package lib

import "strings"

func IsNullOrEmpty(s interface{}) bool {
	return s == nil || strings.TrimSpace(s.(string)) == ""
}
