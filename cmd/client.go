package cmd

import (
	"github.com/MJR0727/plato_mjr/client"
	"github.com/spf13/cobra"
)

// 客户端cmd

func init() {
	rootCmd.AddCommand()
}

var clientCommand = &cobra.Command{
	Use: "client",
	Run: ClientHandler,
}

func ClientHandler(cmd *cobra.Command, args []string) {
	client.RunMain()
}
