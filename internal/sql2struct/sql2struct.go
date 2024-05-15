// Package sql2struct
// @Description: sql生成struct
// @Auth shigx 2024-05-14 17:09:19
package sql2struct

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"tool-cli/internal/mysql"
)

var mysqlTypeToGoType = map[string]string{
	"tinyint":    "int64",
	"smallint":   "int64",
	"mediumint":  "int64",
	"int":        "int64",
	"integer":    "int64",
	"bigint":     "int64",
	"float":      "float64",
	"double":     "float64",
	"decimal":    "float64",
	"date":       "string",
	"time":       "string",
	"year":       "string",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"char":       "string",
	"varchar":    "string",
	"tinyblob":   "string",
	"tinytext":   "string",
	"blob":       "string",
	"text":       "string",
	"mediumblob": "string",
	"mediumtext": "string",
	"longblob":   "string",
	"longtext":   "string",
}

// GetModelTemplate
// @Description 根据模板生成model
// @Auth shigx
// @Date 2022/4/20 6:42 下午
// @param
// @return
func GetModelTemplate(columns []mysql.TableColumn, tableName string, tableComment string) ([]byte, error) {
	t, err := GetTemplate()

	if err != nil {
		return nil, errors.Wrap(err, "template init err")
	}

	var structContent = make([]string, 0)
	for _, row := range columns {
		str := fmt.Sprintf("%s %s %s", Capitalize(row.ColumnName), TextToType(row.DataType), getGormContent(row))
		structContent = append(structContent, str)
	}

	data := map[string]interface{}{
		"pkg":           tableName,
		"structName":    Capitalize(tableName),
		"structComment": tableComment,
		"structContent": structContent,
		"tableName":     tableName,
	}

	buffer := bytes.NewBufferString("")
	err = t.Execute(buffer, data)
	if err != nil {
		return nil, errors.WithMessage(err, "template data err")
	}

	return format.Source(buffer.Bytes())
}

// getGormContent
//
//	@Description: 字段生成gorm信息
//	@Auth shigx 2024-05-15 09:01:09
//	@param row
//	@return string
func getGormContent(row mysql.TableColumn) string {
	str := fmt.Sprintf("`gorm:\"column:%s", row.ColumnName)
	if row.ColumnKey.String == "PRI" {
		str += ";primary_key"
	}
	if row.Extra.String == "auto_increment" {
		str += ";AUTO_INCREMENT"
	}
	if row.IsNullable == "NO" {
		str += ";NOT NULL"
	}

	str += ";default:" + row.ColumnDefault.String + ";comment:'" + strings.ReplaceAll(row.ColumnComment.String, "\n", "") + "'\"`"

	return str
}

// Capitalize
// @Description 带下划线字符串转首字母大写驼峰
// @Auth shigx
// @Date 2022/3/24 10:04 下午
// @param
// @return
func Capitalize(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	titleCaser := cases.Title(language.English)
	s = titleCaser.String(s)
	return strings.ReplaceAll(s, " ", "")
}

// TextToType @Description mysql类型转go结构体类型
// @Auth shigx
// @Date 2022/3/24 10:28 下午
// @param
// @return
func TextToType(s string) string {
	if val, ok := mysqlTypeToGoType[s]; ok {
		return val
	}

	return ""
}
