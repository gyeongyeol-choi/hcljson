package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/tmax-cloud/hcljson/convert"
)

const (
	colorCyan    = "\033[1;36m%s\033[0m"
	colorYellow  = "\033[1;33m%s\033[0m"
	WarningColor = "\033[1;33m[Warning]\033[0m"
	ErrorColor   = "\033[1;31m[Error]\033[0m"
)

var parser = hclparse.NewParser()

// MEMO : 'gopherjs build .' 커맨드 실행으로 js라이브러리화
func main() {
	js.Module.Get("exports").Set("HclToJson", convert.HclToJson)
	js.Module.Get("exports").Set("JsonToHcl", convert.JsonToHcl)
	js.Module.Get("exports").Set("HclToHcl", convert.HclToHcl)
}
