package cmd

import (
	"hello/plato_mjr/ipconfig"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ipconfigCommand)
}

var ipconfigCommand = &cobra.Command{
	Use: "ipconfig",
	Run: IpconfigHandler,
}

func IpconfigHandler(cmd *cobra.Command, args []string) {
	ipconfig.RunMain(ConfigPath)
}
