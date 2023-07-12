package tools

import (
	"strings"
)

// 目前api格式为 {{servername}}/api/{{version}}/{{path}}

func FormatPermissionFromUri(path string, method string) string {
	array := strings.SplitN(path, "/", 5)
	return array[1] + ":" + path + ":" + method
}
