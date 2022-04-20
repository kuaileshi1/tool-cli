// @Title 代码文件注释提取操作
// @Description 根据不同参数提取对应注释并处理
// @Author shigx 2021/10/27 5:26 下午
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"tool-cli/internal/comment"

	"github.com/spf13/cobra"
)

var (
	constType    string // 常量类型
	commentInput string // 输入文件路径
	commentOut   string // 输出文件路径
)

var commentDesc = strings.Join([]string{
	"该命令支持一下模式：",
	"con：常量注释提取，根据常量值获取常量描述信息",
}, "\n")

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "文件注释提取",
	Long:  commentDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

// 常量提取操作
var conCmd = &cobra.Command{
	Use:   "con",
	Short: "提取常量注释并生成map",
	Run: func(cmd *cobra.Command, args []string) {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, viper.GetString("input"), nil, parser.ParseComments)
		checkErr(err)

		var comments = make(map[string]string)
		cmap := ast.NewCommentMap(fset, f, f.Comments)
		for node := range cmap {
			if spec, ok := node.(*ast.ValueSpec); ok && len(spec.Names) == 1 {
				ident := spec.Names[0]
				// 类型为常量处理
				if ident.Obj.Kind == ast.Con {
					// 优先获取单行注释
					switch {
					case spec.Comment != nil:
						comments[ident.Name] = comment.GetConComment(spec.Comment)
					case spec.Doc != nil:
						comments[ident.Name] = comment.GetConComment(spec.Doc)
					}
				}
			}
		}

		pkg := os.Getenv("GOPACKAGE")
		if f.Name != nil {
			pkg = f.Name.Name
		}
		code, err := comment.GetConCode(viper.GetString("type"), pkg, comments)
		checkErr(err)
		if commentOut == "" {
			commentOut = strings.TrimSuffix(viper.GetString("input"), ".go") + "_msg.go"
		}

		checkErr(ioutil.WriteFile(commentOut, code, 0644))
		fmt.Println("处理成功，output:", commentOut)
	},
}

func init() {
	commentCmd.AddCommand(conCmd)

	conCmd.Flags().StringVarP(&commentInput, "input", "i", os.Getenv("GOFILE"), `需要提取的文件`)
	conCmd.Flags().StringVarP(&commentOut, "output", "o", "", `输出文件`)
	conCmd.Flags().StringVarP(&constType, "type", "t", "int", "常量类型")

	viper.BindPFlag("input", conCmd.Flags().Lookup("input"))
	viper.BindPFlag("output", conCmd.Flags().Lookup("output"))
	viper.BindPFlag("type", conCmd.Flags().Lookup("type"))
}
