package tools

import (
	"fmt"
	"regexp"
	"strings"
)

// 目前api格式为 {{servername}}/api/{{version}}/{{path}}:{{method}}

func FormatPermissionFromUri(path string, method string) string {
	//去掉查询参数
	path = strings.Split(path, "?")[0]
	// array := strings.SplitN(path, "/", 3)
	// return array[1] + ":" + path + ":" + method
	return path + ":" + method
}

func MatchRegex(regex string, path string) bool {
	pathRegex, err := regexp.Compile(regex)
	// fmt.Println("pathRegex", pathRegex)
	if err != nil {
		fmt.Println("Invalid regex pattern:", err)
		return false
	}
	// 提取路径参数
	match := pathRegex.FindStringSubmatch(path)
	// fmt.Println("url match result -> ", match)
	if len(match) > 0 {
		for _, v := range match {
			if v == path {
				fmt.Println("match regex success", regex, path)
				return true
			}
		}
	}
	return false
}

type ResourceAuthenticatedItem struct {
	Resource      string `json:"resource"`
	Authenticated bool   `json:"authenticated"`
}
