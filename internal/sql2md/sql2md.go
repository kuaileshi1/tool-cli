// Package sql2md
// @Description: 将sql生成md文件
// @Auth shigx 2024-05-14 17:48:21
package sql2md

import (
	"fmt"
	"strings"
	"tool-cli/internal/mysql"
)

// GetMdContent
// @Description 将表信息生成md格式字符串
// @Auth shigx
// @Date 2024-05-14 17:48:21
// @param
// @return
func GetMdContent(columns []mysql.TableColumn, dbName string, tableName string, tableComment string) string {
	mdContent := fmt.Sprintf("#### %s.%s \n", dbName, tableName)
	if tableComment != "" {
		mdContent += tableComment + "\n"
	}
	mdContent += "\n" +
		"|  序号 |           字段名 |             类型 |    键 |   为空 |                 额外 |      默认值 |                 描述 |\n" +
		"| :---: | :-------------: | :-------------: | :---: | :---: | :------------------: | :--------: | :------------------: |\n"
	for _, row := range columns {
		mdContent += fmt.Sprintf("| %5d | %15s | %15s | %5s | %5s | %20s | %10s | %20s |\n",
			row.OrdinalPosition,
			row.ColumnName,
			row.ColumnType,
			row.ColumnKey.String,
			row.IsNullable,
			row.Extra.String,
			row.ColumnDefault.String,
			strings.ReplaceAll(strings.ReplaceAll(row.ColumnComment.String, "|", "\\|"), "\n", ""),
		)
	}
	return mdContent
}
