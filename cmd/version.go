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
	Short: "Print the version number of masbench",
	Long:  `All software has versions. This is masbench's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("masbench v0.0.1")
	},
}