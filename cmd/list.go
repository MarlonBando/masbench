package cmd

import (
	"fmt"
	"os"

	"masbench/internals/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.Flags().Bool("comparisons", false, "show also comparisons")
}

var listCmd = &cobra.Command{
	Short: "list all benchmarks",
	Use:   "list",
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func list() {
	cfg := config.GetConfig()
	entries, err := os.ReadDir(cfg.BenchmarkFolder)
	if err != nil {
		fmt.Printf("failed to read directory %s: %s\n", cfg.BenchmarkFolder, err.Error())
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if !entry.IsDir() {
			continue
		}

		if entryName == "comparisons" {
			continue
		}

		fmt.Println(entry.Name())
	}

	//TODO: check the flag to include the available comparisons
}
