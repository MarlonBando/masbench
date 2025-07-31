package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"masbench/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize masbench in this repository",
	Long: "This command sets up the necessary configuration files and directories for masbench " +
		"in the current repository.",
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func runInit() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	fmt.Println("masbench initialization will create masbench_config.yml in", currentDir)
	if !askForConfirmation("Do you want to continue? [Y/n]") {
		fmt.Println("Initialization cancelled.")
		return
	}

	prompt := "Do you want to configure the file with an interactive dialog? " +
		"(If not, a default file will be created for manual edit) [Y/n]"
	if askForConfirmation(prompt) {
		interactiveConfigCreation(currentDir)
	} else {
		createDefaultConfig(currentDir)
	}
}

func createDefaultConfig(basePath string) {
	defaultConfig := config.Config{
		ServerPath:      "path/to/your/server_executable",
		LevelsDir:       "path/to/your/levels_directory",
		BenchmarkFolder: "benchmarks",
		ClientCommand:   "your_client_command --level {level_path}",
		Timeout:         "180",
	}

	writeConfig(basePath, &defaultConfig)
	fmt.Printf("Created default config file at %s. Please edit it manually.\n", filepath.Join(basePath, "masbench_config.yml"))
}

func interactiveConfigCreation(basePath string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please provide the following configuration values.")

	serverPath := askForServerPath(reader)
	levelsDir := askForLevelsDir(reader)
	benchmarkFolder := askForBenchmarkFolder(reader)
	clientCommand := askForClientCommand(reader)
	
	newConfig := config.Config{
		ServerPath:      serverPath,
		LevelsDir:       levelsDir,
		BenchmarkFolder: benchmarkFolder,
		ClientCommand:   clientCommand,
		Timeout:         "180", // Default timeout to simplify initialization
	}

	writeConfig(basePath, &newConfig)
	configPath := filepath.Join(basePath, "masbench_config.yml")
	fmt.Printf("Created config file at %s.\n", configPath)
}

func writeConfig(basePath string, cfg *config.Config) {
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshalling config: %v\n", err)
		return
	}

	configPath := filepath.Join(basePath, "masbench_config.yml")
	err = os.WriteFile(configPath, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
	}
}

func askForLevelsDir(reader *bufio.Reader) string {
	prompt := "Enter the path to your level directory. It's recommended to create a separate folder " +
		"where you place only the levels that you want to benchmark to reduce the benchmark time."
	for {
		fmt.Println(prompt)
		fmt.Print("Levels directory path: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Path cannot be empty.")
			continue
		}

		info, err := os.Stat(input)
		if err != nil {
			if os.IsNotExist(err) {
				confirmationPrompt := fmt.Sprintf("Directory '%s' does not exist. Create it? [Y/n]", input)
				if askForConfirmation(confirmationPrompt) {
					if err := os.MkdirAll(input, 0755); err != nil {
						fmt.Printf("Error creating directory '%s': %v\n", input, err)
						continue 
					}
					fmt.Printf("Created directory '%s'.\n", input)
					return input
				} else {
					fmt.Println("Please enter a different path.")
					continue
				}
			} else {
				fmt.Printf("Error checking directory path '%s': %v\n", input, err)
				continue
			}
		}

		if !info.IsDir() {
			fmt.Printf("Invalid input: The path '%s' is not a directory.\n", input)
			continue
		}

		return input
	}
}

func askForServerPath(reader *bufio.Reader) string {
	for {
		fmt.Print("Enter the path to your server executable (.jar file): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if !strings.HasSuffix(strings.ToLower(input), ".jar") {
			fmt.Println("Invalid input: The path must be for a .jar file.")
			continue
		}

		if _, err := os.Stat(input); os.IsNotExist(err) {
			fmt.Printf("Invalid input: The file '%s' does not exist.\n", input)
			continue
		} else if err != nil {
			// Handle other potential errors from os.Stat, though for this case, we can just report and loop.
			fmt.Printf("Error checking file path '%s': %v\n", input, err)
			continue
		}

		return input
	}
}

func askForString(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt + " ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" || input == "" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Invalid input. Please enter y/yes or n/no.")
	}
}

func askForBenchmarkFolder(reader *bufio.Reader) string {
	for {
		fmt.Print("Press Enter to create the benchmark folder in the current directory" +
			" or provide a custom path.")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				return "benchmarks"
			}
			return filepath.Join(currentDir, "benchmarks")
		}

		if _, err := os.Stat(input); os.IsNotExist(err) {
			confirmationPrompt := fmt.Sprintf("Directory '%s' does not exist. Create it? [Y/n]", input)
			if askForConfirmation(confirmationPrompt) {
				if err := os.MkdirAll(input, 0755); err != nil {
					fmt.Printf("Error creating directory '%s': %v\n", input, err)
					continue 
				}
				fmt.Printf("Created directory '%s'.\n", input)
				return input
			} else {
				fmt.Println("Please enter a different path.")
				continue
			}
		} else if err != nil {
			fmt.Printf("Error checking directory path '%s': %v\n", input, err)
			continue
		}

		return input
	}
}

func askForClientCommand(reader *bufio.Reader) string{
	fmt.Print("Enter the command to start your client.\nExample: \"python -m project.src.searchclient -greedy --max-memory 1024\"\n" +
		"Note: This is the command that goes after -c when running the server.\n" +
		"Do not include 'java -jar server.jar' or level path, just your client command: ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}