// Package cmd
// @Description: 将sql建表语句生产struct
// @Auth shigx 2024-05-14 16:35:12
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"tool-cli/internal/mysql"
	"tool-cli/internal/sql2struct"
)

var sql2structCmd = &cobra.Command{
	Use:   "sql2struct",
	Short: "将mysql表生成struct文件",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("mysql.addr", cmd.Flags().Lookup("addr"))
		_ = viper.BindPFlag("mysql.user", cmd.Flags().Lookup("user"))
		_ = viper.BindPFlag("mysql.pass", cmd.Flags().Lookup("pass"))
		_ = viper.BindPFlag("mysql.db", cmd.Flags().Lookup("db"))
		_ = viper.BindPFlag("mysql.table", cmd.Flags().Lookup("table"))
		_ = viper.BindPFlag("mysql.dir", cmd.Flags().Lookup("dir"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := &mysql.Config{
			Addr:     viper.GetString("mysql.addr"),
			User:     viper.GetString("mysql.user"),
			Password: viper.GetString("mysql.pass"),
			DbName:   viper.GetString("mysql.db"),
		}
		db, err := mysql.New(config)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer func() {
			// 关闭数据库连接
			cobra.CheckErr(db.CloseDb())
		}()

		// 检查输出目录是否存在，不存在则创建
		filePath := viper.GetString("mysql.dir")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.MkdirAll(filePath, 0755); err != nil {
				cobra.CheckErr(err)
			}
		}

		// 查询表备注信息
		tableComment, err := mysql.GetTableComment(db.GetDb(), viper.GetString("mysql.db"), viper.GetString("mysql.table"))
		cobra.CheckErr(err)

		// 查询表字段信息
		tableColumn, err := mysql.GetTableColumn(db.GetDb(), viper.GetString("mysql.db"), viper.GetString("mysql.table"))
		cobra.CheckErr(err)

		// 创建model文件
		modelName := path.Join(filePath, viper.GetString("mysql.table")+".go")

		code, err := sql2struct.GetModelTemplate(tableColumn, viper.GetString("mysql.table"), tableComment)
		cobra.CheckErr(err)

		cobra.CheckErr(os.WriteFile(modelName, code, 0644))

		fmt.Println("table:" + viper.GetString("mysql.table") + "生成struct文件完成")
	},
}

func init() {
	var (
		addr, user, password, db, table, out string
	)

	sql2structCmd.Flags().StringVar(&addr, "addr", "127.0.0.1:3306", "请输入db地址，例：127.0.0.1:3306")
	sql2structCmd.Flags().StringVar(&user, "user", "root", "请输入db用户名")
	sql2structCmd.Flags().StringVar(&password, "pass", "", "请输入db密码")
	sql2structCmd.Flags().StringVar(&db, "db", "", "请输入db名称")
	sql2structCmd.Flags().StringVar(&table, "table", "", "请输入表名")
	sql2structCmd.Flags().StringVar(&out, "dir", "./", "请输入输出目录")
}
