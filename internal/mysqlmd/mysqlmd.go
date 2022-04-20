// @Title 请填写文件名称（需要改）
// @Description 请填写文件描述（需要改）
// @Author shigx 2022/3/24 9:22 下午
package mysqlmd

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"go/format"
	"gorm.io/gorm"
	"strings"
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

// @Description 表字段信息定义
// @Auth shigx
// @Date 2022/3/24 10:42 下午
// @param
// @return
type TableColumn struct {
	OrdinalPosition uint16         `gorm:"column:ORDINAL_POSITION"` // 字段顺序
	ColumnName      string         `gorm:"column:COLUMN_NAME"`      // 字段名称
	ColumnType      string         `gorm:"column:COLUMN_TYPE"`      // 字段类型
	DataType        string         `gorm:"column:DATA_TYPE"`        // 数据类型
	ColumnKey       sql.NullString `gorm:"column:COLUMN_KEY"`       // 字段键
	IsNullable      string         `gorm:"column:IS_NULLABLE"`      // 是否允许为空
	Extra           sql.NullString `gorm:"column:EXTRA"`            // 额外信息
	ColumnComment   sql.NullString `gorm:"column:COLUMN_COMMENT"`   // 字段备注
	ColumnDefault   sql.NullString `gorm:"column:COLUMN_DEFAULT"`   // 字段默认值
}

// @Description 查询表备注信息
// @Auth shigx
// @Date 2022/3/24 9:37 下午
// @param
// @return
func GetTableComment(db *gorm.DB, dbName string, tableName string) (string, error) {
	// sql := "SELECT `table_comment` FROM `information_schema`.`tables` WHERE"
	var comment string
	if err := db.Table("information_schema.tables").
		Select("table_comment").
		Where("table_schema = ? and table_name = ?", dbName, tableName).
		Take(&comment).Error; err != nil {
		return "", err
	}

	return comment, nil
}

// @Description 返回表字段信息
// @Auth shigx
// @Date 2022/3/24 10:54 下午
// @param
// @return
func GetTableColumn(db *gorm.DB, dbName string, tableName string) ([]TableColumn, error) {
	ret := make([]TableColumn, 0)
	err := db.Table("information_schema.columns").
		Select(`ORDINAL_POSITION`, `COLUMN_NAME`, `COLUMN_TYPE`, `DATA_TYPE`, `COLUMN_KEY`, `IS_NULLABLE`, `EXTRA`, `COLUMN_COMMENT`, `COLUMN_DEFAULT`).
		Where("table_schema = ? and table_name = ?", dbName, tableName).
		Order("ORDINAL_POSITION ASC").
		Find(&ret).
		Error

	return ret, err
}

// @Description 将表信息生成md格式字符串
// @Auth shigx
// @Date 2022/3/24 11:20 下午
// @param
// @return
func GetMdContent(columns []TableColumn, dbName string, tableName string, tableComment string) string {
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

// @Description 将表生成gorm model
// @Auth shigx
// @Date 2022/4/20 4:26 下午
// @param
// @return
func GetModelContent(columns []TableColumn, tableName string, tableComment string) ([]byte, error) {
	packageContent := fmt.Sprintf("package %s\n", tableName)
	structContent := fmt.Sprintf("\n\n// %s %s \n", Capitalize(tableName), tableComment)
	structContent += fmt.Sprintf("type %s struct {\n", Capitalize(tableName))

	for _, row := range columns {
		structContent += fmt.Sprintf("%s %s `gorm:\"%s\"` // %s\n", Capitalize(row.ColumnName), TextToType(row.DataType), row.ColumnName, strings.ReplaceAll(row.ColumnComment.String, "\n", ""))
	}
	structContent += "}\n"
	structContent += fmt.Sprintf("func (%s) TableName() string {\n", Capitalize(tableName))
	structContent += fmt.Sprintf("return \"%s\"", tableName)
	structContent += "}\n"

	// 如果有引入 time.Time, 则需要引入 time 包
	var importContent string
	if strings.Contains(structContent, "time.Time") {
		importContent = "import \"time\"\n\n"
	}

	return format.Source([]byte(packageContent + importContent + structContent))
}

// @Description 根据模板生成model（方法二）
// @Auth shigx
// @Date 2022/4/20 6:42 下午
// @param
// @return
func GetModelTemplate(columns []TableColumn, tableName string, tableComment string) ([]byte, error) {
	t, err := GetTemplate()

	if err != nil {
		return nil, errors.Wrap(err, "template init err")
	}

	var structContent = make([]string, 0)
	for _, row := range columns {
		str := fmt.Sprintf("%s %s `gorm:\"%s\"` // %s", Capitalize(row.ColumnName), TextToType(row.DataType), row.ColumnName, strings.ReplaceAll(row.ColumnComment.String, "\n", ""))
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

// @Description 带下划线字符串转首字母大写驼峰
// @Auth shigx
// @Date 2022/3/24 10:04 下午
// @param
// @return
func Capitalize(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(s)
	return strings.ReplaceAll(s, " ", "")
}

// @Description mysql类型转go结构体类型
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
