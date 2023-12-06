package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// 每条命令初始化执行配置。
	cobra.OnInitialize(initConfig)
}

func initConfig() {

}

var rootCmd = &cobra.Command{
	Use:   "plato_mjr",
	Short: "一个值得学习的IM系统",
	Run:   plato,
}

// 命令行初始化函数
func plato(cmd *cobra.Command, args []string) {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
