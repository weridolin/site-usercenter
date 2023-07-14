package tools

import (
	"strings"
)

// 目前api格式为 {{servername}}/api/{{version}}/{{path}}

func FormatPermissionFromUri(path string, method string) string {
	//去掉查询参数
	path = strings.Split(path, "?")[0]
	array := strings.SplitN(path, "/", 3)
	return array[1] + ":" + path + ":" + method
}
