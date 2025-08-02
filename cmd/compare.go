package cmd

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"masbench/config"
	"github.com/fogleman/gg"
	"github.com/go-gota/gota/dataframe"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func init() {
	rootCmd.AddCommand(compareCmd)
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two benchmark results",
	Long:  `This command compares two benchmark results and create tables and graphs with the comparison.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(colorRed + "Error: You must provide two benchmark result files to compare." + colorReset)
			os.Exit(1)
		}
	
		benchmark1Name := args[0]
		Benchmark2Name := args[1]
		
		cfg := config.GetConfig()

		benchmark1Path := filepath.Join(cfg.BenchmarkFolder, benchmark1Name, fmt.Sprintf("%s_results.csv", benchmark1Name))
		benchmark2Path := filepath.Join(cfg.BenchmarkFolder, Benchmark2Name, fmt.Sprintf("%s_results.csv", Benchmark2Name))

		if _, err := os.Stat(benchmark1Path); os.IsNotExist(err) {
			fmt.Printf(colorRed+"Error: Benchmark result file not found: %s%s", benchmark1Name, colorReset)
			os.Exit(1)
		}

		if _, err := os.Stat(benchmark2Path); os.IsNotExist(err) {
			fmt.Printf(colorRed+"Error: Benchmark result file not found: %s%s", Benchmark2Name, colorReset)
			os.Exit(1)
		}

		compareResults(benchmark1Path, benchmark2Path, benchmark1Name, Benchmark2Name)
	},
}

func compareResults(benchmark1Path, benchmark2Path, name1, name2 string) {
	df1, err := readCSV(benchmark1Path)
	if err != nil {
		fmt.Printf(colorRed+"Error reading benchmark1 CSV: %v%s\n", err, colorReset)
		os.Exit(1)
	}

	df2, err := readCSV(benchmark2Path)
	if err != nil {
		fmt.Printf(colorRed+"Error reading benchmark2 CSV: %v%s\n", err, colorReset)
		os.Exit(1)
	}

	cfg := config.GetConfig()
	outputDir := filepath.Join(cfg.BenchmarkFolder, "comparisons", fmt.Sprintf("%svs%s", name1, name2))
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf(colorRed+"Error creating output directory: %v%s\n", err, colorReset)
		os.Exit(1)
	}

	// Create bar charts for each metric
	metrics := []string{"Generated", "Explored", "MemoryAlloc", "Time", "Actions"}
	for _, metric := range metrics {
		chartPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", metric))
		if err := createBarChart(df1, df2, metric, name1, name2, chartPath); err != nil {
			fmt.Printf(colorRed+"Error creating %s chart: %v%s\n", metric, err, colorReset)
		} else {
			fmt.Printf(colorGreen+"Created chart: %s%s\n", chartPath, colorReset)
		}
	}

	// Generate difference report table
	tablePath := filepath.Join(outputDir, fmt.Sprintf("%svs%s_table.png", name1, name2))
	if err := generateDifferenceReport(df1, df2, name1, name2, tablePath); err != nil {
		fmt.Printf(colorRed+"Error creating difference report: %v%s\n", err, colorReset)
	} else {
		fmt.Printf(colorGreen+"Created difference report: %s%s\n", tablePath, colorReset)
	}

	fmt.Printf(colorGreen+"Comparison completed successfully! Results saved to: %s%s\n", outputDir, colorReset)
}

func readCSV(filePath string) (dataframe.DataFrame, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	defer file.Close()

	return dataframe.ReadCSV(file), nil
}

func createBarChart(df1, df2 dataframe.DataFrame, metric, name1, name2, outputPath string) error {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Comparison of %s", metric)
	p.Y.Label.Text = metric
	p.X.Label.Text = "Level Name"

	// Extract data for the specified metric
	levelNames1 := df1.Col("LevelName").Records()
	metricValues1 := df1.Col(metric).Records()
	
	// Create a map for df2 data lookup
	df2Map := make(map[string]string)
	levelNames2 := df2.Col("LevelName").Records()
	metricValues2 := df2.Col(metric).Records()
	
	for i, level := range levelNames2 {
		if i < len(metricValues2) {
			df2Map[level] = metricValues2[i]
		}
	}
	
	// Prepare data for plotting
	var levels []string
	var values1 []float64
	var values2 []float64
	
	for i, level := range levelNames1 {
		if i < len(metricValues1) {
			levels = append(levels, level)
			
			// Parse value from df1
			val1, err1 := strconv.ParseFloat(metricValues1[i], 64)
			if err1 != nil {
				val1 = 0
			}
			values1 = append(values1, val1)
			
			// Parse value from df2
			val2 := 0.0
			if val2Str, exists := df2Map[level]; exists {
				if parsedVal2, err2 := strconv.ParseFloat(val2Str, 64); err2 == nil {
					val2 = parsedVal2
				}
			}
			values2 = append(values2, val2)
		}
	}
	
	// Create bar chart data
	bars1, err := plotter.NewBarChart(plotter.Values(values1), vg.Points(20))
	if err != nil {
		return err
	}
	bars1.Color = plotutil.Color(0) // Blue
	bars1.Offset = -vg.Points(10)
	
	bars2, err := plotter.NewBarChart(plotter.Values(values2), vg.Points(20))
	if err != nil {
		return err
	}
	bars2.Color = plotutil.Color(1) // Orange
	bars2.Offset = vg.Points(10)
	
	p.Add(bars1, bars2)
	p.Legend.Add(name1, bars1)
	p.Legend.Add(name2, bars2)
	
	// Set X-axis labels
	p.NominalX(levels...)
	
	// Save the plot
	if err := p.Save(12*vg.Inch, 6*vg.Inch, outputPath); err != nil {
		return err
	}
	
	return nil
}

func generateDifferenceReport(df1, df2 dataframe.DataFrame, name1, name2, outputPath string) error {
	// Get the canvas dimensions
	const width = 1200
	const height = 800
	
	// Create a new context
	dc := gg.NewContext(width, height)
	
	// Set background to white
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.Clear()
	
	// Set up drawing parameters
	dc.SetColor(color.RGBA{0, 0, 0, 255}) // Black text
	fontSize := 12.0
	dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", fontSize)
	
	// Title
	title := fmt.Sprintf("Benchmark Comparison: %s vs %s", name1, name2)
	dc.DrawStringAnchored(title, width/2, 50, 0.5, 0.5)
	
	// Table headers
	headers := []string{"LevelName", "Generated (Diff)", "Explored (Diff)", "MemoryAlloc (Diff)", "Time (Diff)", "Actions (Diff)", "Solved"}
	headerY := 100
	colWidth := width / len(headers)
	
	// Draw headers
	for i, header := range headers {
		x := float64(i * colWidth + colWidth/2)
		dc.DrawStringAnchored(header, x, float64(headerY), 0.5, 0.5)
	}
	
	// Draw header line
	dc.DrawLine(0, float64(headerY+20), float64(width), float64(headerY+20))
	dc.Stroke()
	
	// Get data from dataframes
	levelNames1 := df1.Col("LevelName").Records()
	
	// Process each row
	rowHeight := 25
	currentY := headerY + 40
	
	// Create a map for easier lookup of df2 data
	df2Map := make(map[string]map[string]string)
	for i := 0; i < df2.Nrow(); i++ {
		levelName := df2.Elem(i, 0).String() // First column is LevelName
		df2Map[levelName] = make(map[string]string)
		for j, colName := range df2.Names() {
			df2Map[levelName][colName] = df2.Elem(i, j).String()
		}
	}
	
	// Process each row from df1
	for i := 0; i < df1.Nrow(); i++ {
		levelName := levelNames1[i]
		y := float64(currentY + i*rowHeight)
		
		// Draw level name
		dc.SetColor(color.RGBA{0, 0, 0, 255})
		dc.DrawStringAnchored(levelName, float64(colWidth/2), y, 0.5, 0.5)
		
		// Process each numeric column
		numericCols := []string{"Generated", "Explored", "MemoryAlloc", "Time", "Actions"}
		for j, col := range numericCols {
			x := float64((j+1)*colWidth + colWidth/2)
			
			// Get column index for df1
			colIndex := getColumnIndex(df1, col)
			if colIndex == -1 {
				continue
			}
			
			val1Str := df1.Elem(i, colIndex).String()
			val1, err1 := strconv.ParseFloat(val1Str, 64)
			
			var diff float64
			var pct float64
			var diffText string
			var bgColor color.RGBA
			
			if df2Data, exists := df2Map[levelName]; exists {
				val2Str := df2Data[col]
				val2, err2 := strconv.ParseFloat(val2Str, 64)
				
				if err1 == nil && err2 == nil {
					diff = val1 - val2
					if val1 != 0 {
						pct = (diff / val1) * 100
					}
					diffText = fmt.Sprintf("%.2f (%.2f%%)", diff, pct)
					
					// Set background color based on difference
					if diff > 0 {
						bgColor = color.RGBA{255, 200, 200, 255} // Light red
					} else if diff < 0 {
						bgColor = color.RGBA{200, 255, 200, 255} // Light green
					} else {
						bgColor = color.RGBA{255, 255, 255, 255} // White
					}
				} else {
					diffText = "N/A"
					bgColor = color.RGBA{255, 255, 255, 255}
				}
			} else {
				diffText = "Missing"
				bgColor = color.RGBA{255, 255, 200, 255} // Light yellow
			}
			
			// Draw background color
			dc.SetColor(bgColor)
			dc.DrawRectangle(float64(j+1)*float64(colWidth), y-float64(rowHeight/2), float64(colWidth), float64(rowHeight))
			dc.Fill()
			
			// Draw text
			dc.SetColor(color.RGBA{0, 0, 0, 255})
			dc.DrawStringAnchored(diffText, x, y, 0.5, 0.5)
		}
		
		// Handle Solved column
		x := float64(6*colWidth + colWidth/2)
		solvedColIndex := getColumnIndex(df1, "Solved")
		var solved1 string
		if solvedColIndex != -1 {
			solved1 = df1.Elem(i, solvedColIndex).String()
		}
		var solvedText string
		var bgColor color.RGBA
		
		if df2Data, exists := df2Map[levelName]; exists {
			solved2 := df2Data["Solved"]
			solvedText = fmt.Sprintf("%s <- %s", solved1, solved2)
			
			// Set background color based on solved transition
			if solved1 == "Yes" && solved2 == "No" {
				bgColor = color.RGBA{200, 255, 200, 255} // Light green
			} else if solved1 == "No" && solved2 == "Yes" {
				bgColor = color.RGBA{255, 200, 200, 255} // Light red
			} else {
				bgColor = color.RGBA{255, 255, 255, 255} // White
			}
		} else {
			solvedText = fmt.Sprintf("%s <- Missing", solved1)
			bgColor = color.RGBA{255, 255, 200, 255} // Light yellow
		}
		
		// Draw background color
		dc.SetColor(bgColor)
		dc.DrawRectangle(float64(6*colWidth), y-float64(rowHeight/2), float64(colWidth), float64(rowHeight))
		dc.Fill()
		
		// Draw text
		dc.SetColor(color.RGBA{0, 0, 0, 255})
		dc.DrawStringAnchored(solvedText, x, y, 0.5, 0.5)
		
		// Draw row separator line
		dc.SetColor(color.RGBA{200, 200, 200, 255})
		dc.DrawLine(0, y+float64(rowHeight/2), float64(width), y+float64(rowHeight/2))
		dc.Stroke()
	}
	
	// Draw vertical lines for columns
	dc.SetColor(color.RGBA{200, 200, 200, 255})
	for i := 1; i < len(headers); i++ {
		x := float64(i * colWidth)
		dc.DrawLine(x, float64(headerY), x, float64(currentY+df1.Nrow()*rowHeight))
		dc.Stroke()
	}
	
	// Save the image
	return dc.SavePNG(outputPath)
}

func getColumnIndex(df dataframe.DataFrame, colName string) int {
	for i, name := range df.Names() {
		if name == colName {
			return i
		}
	}
	return -1
}
