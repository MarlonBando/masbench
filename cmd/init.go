package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use: "init",
	Short: "Initialize masbench in this repository",
	Long: `This command sets up the necessary configuration files and directories for masbench in the current repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("masbench initialized")
	},
}