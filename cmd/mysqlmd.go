// @Title 将mysql数据表生成struct和md文档
// @Description 将mysql数据表生成struct和md文档
// @Author shigx 2022/3/24 2:35 下午
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"tool-cli/internal/mysql"
	"tool-cli/internal/mysqlmd"
)

// 数据库信息定义
var (
	dbAddr  string
	dbUser  string
	dbPass  string
	dbName  string
	dbTable string
	dir     string
)

var mysqlmdCmd = &cobra.Command{
	Use:   "mysqlmd",
	Short: "将mysql表生成md及struct文件",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := mysql.New(viper.GetString("mysqlmd.addr"),
			viper.GetString("mysqlmd.user"),
			viper.GetString("mysqlmd.pass"),
			viper.GetString("mysqlmd.name"))
		if err != nil {
			checkErr(err)
		}
		defer func() {
			// 关闭数据库连接
			err := db.CloseDb()
			checkErr(err)
		}()

		// 检查输出目录是否存在，不存在则创建
		filePath := strings.TrimRight(viper.GetString("mysqlmd.dir"), "/") + "/" + dbTable
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			os.MkdirAll(filePath, 0755)
		}

		// 创建md文件
		mdFileName := path.Join(filePath, "gen_table.md")
		mdFile, err := os.OpenFile(mdFileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
		if err != nil {
			checkErr(err)
		}

		// 查询表备注信息
		tableComment, err := mysqlmd.GetTableComment(db.GetDb(), viper.GetString("mysqlmd.name"), dbTable)
		checkErr(err)

		tableColumn, err := mysqlmd.GetTableColumn(db.GetDb(), viper.GetString("mysqlmd.name"), dbTable)
		checkErr(err)
		mdContent := mysqlmd.GetMdContent(tableColumn, viper.GetString("mysqlmd.name"), dbTable, tableComment)

		mdFile.WriteString(mdContent)
		mdFile.Close()
		fmt.Println("table:" + dbTable + "生成md文件完成")

		// 创建model文件
		modelName := path.Join(filePath, "gen_model.go")
		// code, err := mysqlmd.GetModelContent(tableColumn, dbTable, tableComment)

		code, err := mysqlmd.GetModelTemplate(tableColumn, dbTable, tableComment)
		checkErr(err)
		checkErr(ioutil.WriteFile(modelName, code, 0644))

		fmt.Println("table:" + dbTable + "生成model文件完成")

	},
}

func init() {
	mysqlmdCmd.Flags().StringVar(&dbAddr, "addr", "127.0.0.1:3306", "请输入db地址，例：127.0.0.1:3306")
	mysqlmdCmd.Flags().StringVar(&dbUser, "user", "root", "请输入db用户名")
	mysqlmdCmd.Flags().StringVar(&dbPass, "pass", "", "请输入db密码")
	mysqlmdCmd.Flags().StringVar(&dbName, "name", "", "请输入db名称")
	mysqlmdCmd.Flags().StringVar(&dbTable, "table", "", "请输入表名")
	mysqlmdCmd.Flags().StringVar(&dir, "dir", "./", "请输入输出目录")

	viper.BindPFlag("mysqlmd.addr", mysqlmdCmd.Flags().Lookup("addr"))
	viper.BindPFlag("mysqlmd.user", mysqlmdCmd.Flags().Lookup("user"))
	viper.BindPFlag("mysqlmd.pass", mysqlmdCmd.Flags().Lookup("pass"))
	viper.BindPFlag("mysqlmd.name", mysqlmdCmd.Flags().Lookup("name"))
	viper.BindPFlag("mysqlmd.table", mysqlmdCmd.Flags().Lookup("table"))
	viper.BindPFlag("mysqlmd.dir", mysqlmdCmd.Flags().Lookup("dir"))
}
