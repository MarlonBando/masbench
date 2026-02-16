package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"masbench/internals/config"
	"masbench/internals/parsers"

	"github.com/spf13/cobra"
)

var message string
var algorithm string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&message, "message", "m", "", "Add a note to the run")
	runCmd.Flags().StringVarP(&algorithm, "algorithm", "a", "", "Algorithm to use for this run")
}

var runCmd = &cobra.Command{
	Use:   "run [benchmark-name]",
	Short: "Run a benchmark with masbench",
	Long: `Execute a benchmark using masbench. Requires a configuration file (masbench_config.yml) to be present in the current directory.

DESCRIPTION
       The run command executes benchmarks against your client implementation.
       It spawns a server process that runs your client against a set of
       test levels, collecting performance metrics and logs.

       Before running, ensure your configuration file is properly set up with
       the server path, levels directory, and client command.

OPTIONS
       -a <algorithm>, --algorithm=<algorithm>
           Specify the algorithm to use for this run (e.g., bfs, dfs, greedy,
           astar). The algorithm argument will be appended to your client
           command using the format defined by AlgorithmFlagFormat in the
           configuration file.

           IMPORTANT: Remove any algorithm flags from your ClientCommand in
           the configuration file before using this option to avoid conflicts.

           EXAMPLE
               If your ClientCommand is:
                   "python -m searchclient.searchclient -bfs"
               Remove "-bfs" and use:
                   masbench run my-benchmark -a bfs

       -m <message>, --message=<message>
           Add a descriptive note or comment to the benchmark run. This
           message will be saved in the benchmark results for reference.

           Useful for documenting the purpose of a run, configuration
           changes, or any other relevant information about the benchmark.

EXAMPLES
       Run a benchmark named "test-run":
           masbench run test-run

       Run with a specific algorithm:
           masbench run astar-test -a astar

       Run with a descriptive message:
           masbench run baseline -m "Baseline performance test"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		benchmarkName := args[0]
		fmt.Printf("Running benchmark: %s\n", benchmarkName)
		runBenchmark(benchmarkName, message, algorithm)
	},
}

func runBenchmark(name, message, algorithm string) {
	cfg := config.GetConfig()

	// Create benchmark folder if it does not exist
	if _, err := os.Stat(cfg.BenchmarkFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.BenchmarkFolder, os.ModePerm); err != nil {
			fmt.Printf("\033[31mError creating benchmark folder: %v\033[0m\n", err)
			return
		}
	}

	benchmarkPath := filepath.Join(cfg.BenchmarkFolder, name)

	// If the benchmark folder for that name already exists, print an error and exit.
	if _, err := os.Stat(benchmarkPath); !os.IsNotExist(err) {
		fmt.Printf("\033[31mError: Benchmark with name '%s' already exists. Please remove it before running a new one.\033[0m\n", name)
		return
	}

	logServerPath := filepath.Join(benchmarkPath, "logs", fmt.Sprintf("%s_server.zip", name))
	logClientPath := filepath.Join(benchmarkPath, "logs", fmt.Sprintf("%s_client.clog", name))

	// Create the logs directory
	logDir := filepath.Dir(logServerPath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("\033[31mError creating log directory: %v\033[0m\n", err)
		return
	}

	// Create the client log file
	logFile, err := os.Create(logClientPath)
	if err != nil {
		fmt.Printf("\033[31mError creating client log file: %v\033[0m\n", err)
		return
	}
	defer logFile.Close()

	if algorithm != "" {
		if strings.Count(cfg.AlgorithmFlagFormat, "%s") != 1 {
			fmt.Println("\033[31mError in your configuration: The parameter AlgorithmFlagFormat in your masbench_config.yml must contain only one %s'\033[0m")
			return
		}
		cfg.ClientCommand += " " + fmt.Sprintf(cfg.AlgorithmFlagFormat, algorithm)
	}

	cmd := exec.Command("java", "-jar", cfg.ServerPath,
		"-l", cfg.LevelsDir,
		"-o", logServerPath,
		"-c", cfg.ClientCommand,
		"-t", fmt.Sprintf("%d", cfg.Timeout),
	)

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	err = cmd.Run()
	if err != nil {
		fmt.Printf("\033[31mError running benchmark: %v\033[0m\n", err)
		os.RemoveAll(benchmarkPath)
		return
	}

	fmt.Println("\033[32mBenchmark run completed successfully.\033[0m")

	// After the benchmark, parse the client log to CSV
	csvOutputPath := filepath.Join(benchmarkPath, fmt.Sprintf("%s_results.csv", name))
	err = parsers.ParseLogToCSV(logClientPath, csvOutputPath)
	if err != nil {
		fmt.Printf("\033[31mError parsing log to CSV: %v\033[0m\n", err)
		return
	}

	descriptionFilePath := filepath.Join(benchmarkPath, name+".md")
	err = os.WriteFile(descriptionFilePath, []byte(message+"\n"), 0644)
	if err != nil {
		fmt.Printf("\033[31mError! Couldn't write in %s \n %v\033[0m\n", descriptionFilePath, err)
	}

	fmt.Printf("\033[32mResults successfully written to %s\033[0m\n", csvOutputPath)
}
