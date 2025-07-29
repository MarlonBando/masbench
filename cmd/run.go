package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var message string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&message, "message", "m", "", "Add a note to the run")
}

var runCmd = &cobra.Command{
	Use:   "run [benchmark-name]",
	Short: "Run a benchmark with masbench",
	Long:  `This command executes a benchmark using masbench. It requires a configuration file to be present in the current directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		benchmarkName := args[0]
		fmt.Printf("Running benchmark: %s\n", benchmarkName)
		if message != "" {
			fmt.Printf("Message: %s\n", message)
		}
		// Here you would typically call a function to run the benchmark
		// For example: runBenchmark(benchmarkName)
	},
}
