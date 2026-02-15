package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"masbench/internals/config"
	"masbench/internals/parsers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refreshCmd)
	refreshCmd.Flags().BoolP("all", "A", false, "Refresh all benchmarks")
}

var refreshCmd = &cobra.Command{
	Use:   "refresh [benchmark-name]",
	Short: "Regenerate CSV results from an existing benchmark's client log",
	Long: `Re-parse the client log file for an existing benchmark and regenerate
the CSV results file. This is useful when the log parser has been updated
and you want to refresh the CSV without re-running the benchmark.

Use --all to refresh all benchmarks at once.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			refreshAll()
		} else if len(args) == 1 {
			refreshBenchmark(args[0])
		} else {
			fmt.Println("\033[31mError: Please specify a benchmark name or use --all.\033[0m")
		}
	},
}

func refreshAll() {
	cfg := config.GetConfig()
	entries, err := os.ReadDir(cfg.BenchmarkFolder)
	if err != nil {
		fmt.Printf("\033[31mError reading benchmark folder: %v\033[0m\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "comparisons" || entry.Name() == "summaries" {
			continue
		}
		refreshBenchmark(entry.Name())
	}
}

func refreshBenchmark(name string) {
	cfg := config.GetConfig()

	benchmarkPath := filepath.Join(cfg.BenchmarkFolder, name)
	if _, err := os.Stat(benchmarkPath); os.IsNotExist(err) {
		fmt.Printf("\033[31mError: Benchmark '%s' not found.\033[0m\n", name)
		return
	}

	logClientPath := filepath.Join(benchmarkPath, "logs", fmt.Sprintf("%s_client.clog", name))
	if _, err := os.Stat(logClientPath); os.IsNotExist(err) {
		fmt.Printf("\033[31mError: Client log file not found: %s\033[0m\n", logClientPath)
		return
	}

	csvOutputPath := filepath.Join(benchmarkPath, fmt.Sprintf("%s_results.csv", name))
	err := parsers.ParseLogToCSV(logClientPath, csvOutputPath)
	if err != nil {
		fmt.Printf("\033[31mError parsing log to CSV: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[32mResults successfully regenerated: %s\033[0m\n", csvOutputPath)
}
