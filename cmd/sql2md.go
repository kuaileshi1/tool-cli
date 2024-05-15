// Package cmd
// @Title 将mysql数据表生成struct和md文档
// @Description 将mysql数据表生成struct和md文档
// @Author shigx 2022/3/24 2:35 下午
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"tool-cli/internal/mysql"
	"tool-cli/internal/sql2md"
)

var sql2mdCmd = &cobra.Command{
	Use:   "sql2md",
	Short: "将mysql表生成md文件",
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
			cobra.CheckErr(os.MkdirAll(filePath, 0755))
		}

		// 查询表备注信息
		tableComment, err := mysql.GetTableComment(db.GetDb(), viper.GetString("mysql.db"), viper.GetString("mysql.table"))
		cobra.CheckErr(err)

		tableColumn, err := mysql.GetTableColumn(db.GetDb(), viper.GetString("mysql.db"), viper.GetString("mysql.table"))
		cobra.CheckErr(err)
		mdContent := sql2md.GetMdContent(tableColumn, viper.GetString("mysql.db"), viper.GetString("mysql.table"), tableComment)

		// 创建md文件
		mdFileName := path.Join(filePath, viper.GetString("mysql.table")+".md")
		err = os.WriteFile(mdFileName, []byte(mdContent), 0644)
		if err != nil {
			cobra.CheckErr(err)
		}

		fmt.Println("table:" + viper.GetString("mysql.table") + "生成md文件完成")
	},
}

func init() {
	var (
		addr, user, pass, db, table, dir string
	)
	sql2mdCmd.Flags().StringVar(&addr, "addr", "127.0.0.1:3306", "请输入db地址，例：127.0.0.1:3306")
	sql2mdCmd.Flags().StringVar(&user, "user", "root", "请输入db用户名")
	sql2mdCmd.Flags().StringVar(&pass, "pass", "", "请输入db密码")
	sql2mdCmd.Flags().StringVar(&db, "db", "", "请输入db名称")
	sql2mdCmd.Flags().StringVar(&table, "table", "", "请输入表名")
	sql2mdCmd.Flags().StringVar(&dir, "dir", "./", "请输入输出目录")
}
