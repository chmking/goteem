package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the Horde CLI and Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Horde load tester -- HEAD")
	},
}
