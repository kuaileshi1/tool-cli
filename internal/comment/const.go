// @Title 常量注释处理
// @Description 将解析的注释信息映射到模版并输出
// @Author shigx 2021/10/28 9:40 上午
package comment

import (
	"bytes"
	"github.com/pkg/errors"
	"go/ast"
	"go/format"
	"strings"
	"text/template"
)

const tpl = `// Code generated by tool-cli DO NOT EDIT
// Package {{.pkg}} const code comment msg
package {{.pkg}}
// noMsg if code is not found, GetMsg will return this
const noMsg = "unknown"
// messages get msg from const comment
var messages = map[{{.constType}}]string{
	{{range $key, $value := .comments}}
	{{$key}}: "{{$value}}",{{end}}
}
{{ if ne .constType "int" }}
// String return string
func (code {{.constType}}) String () string {
	return GetMsg(code)
}
{{ end }}
// GetMsg get error msg
func GetMsg(code {{.constType}}) string {
	var (
		msg string
		ok  bool
	)
	if msg, ok = messages[code]; !ok {
		msg = noMsg
	}
	return msg
}`

// @Description 注释处理
// @Auth shigx
// @Date 2021/10/28 9:58 上午
// @param
// @return
func GetConComment(group *ast.CommentGroup) string {
	var buf bytes.Buffer
	for key, comment := range group.List {
		text := strings.TrimPrefix(comment.Text, "//")
		if key == 0 {
			text = strings.TrimSpace(text)
		}
		buf.WriteString(text)
	}

	return buf.String()
}

// @Description 将数据填充模版并返回
// @Auth shigx
// @Date 2021/10/28 10:24 上午
// @param constType string 常量类型
// @param pkg string 包名
// @param comments 常量注释信息
// @return
func GetConCode(constType string, pkg string, comments map[string]string) ([]byte, error) {
	var (
		t   *template.Template
		err error
		buf = bytes.NewBufferString("")
	)
	data := map[string]interface{}{
		"pkg":       pkg,
		"comments":  comments,
		"constType": constType,
	}
	t, err = template.New("").Parse(tpl)
	if err != nil {
		return nil, errors.Wrap(err, "template init err")
	}
	err = t.Execute(buf, data)
	if err != nil {
		return nil, errors.WithMessage(err, "template data err")
	}

	return format.Source(buf.Bytes())
}
