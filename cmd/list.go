package cmd

import (
	"os"
)

var listCmd = &cobra.Command{
	Use:   "masbench list",
	Short: "List all the avilable benchamarks",
	Long:  `List all the avilable benchamarks`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func Execute() {
}
