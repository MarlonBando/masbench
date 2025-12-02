package comparator

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/go-gota/gota/dataframe"
)

func GenerateHTMLReport(df1, df2 dataframe.DataFrame, name1, name2, outputPath string) error {
	report := prepareComparisonData(df1, df2, name1, name2)

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, report); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func prepareComparisonData(df1, df2 dataframe.DataFrame, name1, name2 string) ComparisonReport {
	report := ComparisonReport{
		Title:          "Benchmark Comparison Report",
		Benchmark1Name: name1,
		Benchmark2Name: name2,
		GeneratedAt:    time.Now().Format("2006-01-02 15:04:05"),
		MetricNames:    []string{"Generated", "Explored", "MemoryAlloc", "Time", "Actions"},
	}

	df2Map := make(map[string]map[string]string)
	for i := 0; i < df2.Nrow(); i++ {
		levelName := df2.Elem(i, 0).String()
		rowData := make(map[string]string)

		for j, colName := range df2.Names() {
			rowData[colName] = df2.Elem(i, j).String()
		}

		df2Map[levelName] = rowData
	}

	levelNames1 := df1.Col("LevelName").Records()

	for i := 0; i < df1.Nrow(); i++ {
		levelName := levelNames1[i]
		levelComp := LevelComparison{
			LevelName: levelName,
		}

		df2Data, exists := df2Map[levelName]

		levelComp.Generated = compareMetric(
			getFloatValue(df1, i, "Generated"),
			getFloatValueFromMap(df2Data, "Generated", exists),
			true,
		)

		levelComp.Explored = compareMetric(
			getFloatValue(df1, i, "Explored"),
			getFloatValueFromMap(df2Data, "Explored", exists),
			true,
		)

		levelComp.MemoryAlloc = compareMetric(
			getFloatValue(df1, i, "MemoryAlloc"),
			getFloatValueFromMap(df2Data, "MemoryAlloc", exists),
			true,
		)

		levelComp.Time = compareMetric(
			getFloatValue(df1, i, "Time"),
			getFloatValueFromMap(df2Data, "Time", exists),
			true,
		)

		levelComp.Actions = compareMetric(
			getFloatValue(df1, i, "Actions"),
			getFloatValueFromMap(df2Data, "Actions", exists),
			true,
		)

		solved1 := getStringValue(df1, i, "Solved")
		solved2 := getStringValueFromMap(df2Data, "Solved", exists)
		levelComp.Solved = compareSolved(solved1, solved2)

		if solved1 == "Yes" && solved2 == "No" {
			levelComp.Generated.Status = "improvement"
			levelComp.Generated.IsImprovement = true
			levelComp.Explored.Status = "improvement"
			levelComp.Explored.IsImprovement = true
			levelComp.MemoryAlloc.Status = "improvement"
			levelComp.MemoryAlloc.IsImprovement = true
			levelComp.Time.Status = "improvement"
			levelComp.Time.IsImprovement = true
			levelComp.Actions.Status = "improvement"
			levelComp.Actions.IsImprovement = true
		} else if solved1 == "No" && solved2 == "Yes" {
			levelComp.Generated.Status = "regression"
			levelComp.Generated.IsImprovement = false
			levelComp.Explored.Status = "regression"
			levelComp.Explored.IsImprovement = false
			levelComp.MemoryAlloc.Status = "regression"
			levelComp.MemoryAlloc.IsImprovement = false
			levelComp.Time.Status = "regression"
			levelComp.Time.IsImprovement = false
			levelComp.Actions.Status = "regression"
			levelComp.Actions.IsImprovement = false
		}

		report.Levels = append(report.Levels, levelComp)
	}

	return report
}

func compareMetric(val1, val2 float64, lowerIsBetter bool) MetricComparison {
	diff := val1 - val2
	var diffPct float64
	if val2 != 0 {
		diffPct = (diff / val2) * 100
	}

	var status string
	var isImprovement bool

	if diff == 0 {
		status = "unchanged"
		isImprovement = false
	} else if lowerIsBetter {
		if diff < 0 {
			status = "improvement"
			isImprovement = true
		} else {
			status = "regression"
			isImprovement = false
		}
	}

	return MetricComparison{
		Value1:        val1,
		Value2:        val2,
		Diff:          diff,
		DiffPct:       diffPct,
		Status:        status,
		IsImprovement: isImprovement,
	}
}

func compareSolved(solved1, solved2 string) SolvedComparison {
	changed := solved1 != solved2
	var status string

	if !changed {
		status = "unchanged"
	} else if solved1 == "Yes" && solved2 == "No" {
		status = "improvement"
	} else {
		status = "regression"
	}

	return SolvedComparison{
		Solved1: solved1,
		Solved2: solved2,
		Changed: changed,
		Status:  status,
	}
}

func getFloatValue(df dataframe.DataFrame, row int, colName string) float64 {
	colIndex := getColumnIndex(df, colName)
	if colIndex == -1 {
		return 0.0
	}

	valStr := df.Elem(row, colIndex).String()
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0.0
	}

	return val
}

func getFloatValueFromMap(data map[string]string, colName string, exists bool) float64 {
	if !exists {
		return 0.0
	}

	valStr, ok := data[colName]
	if !ok {
		return 0.0
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0.0
	}

	return val
}

func getStringValue(df dataframe.DataFrame, row int, colName string) string {
	colIndex := getColumnIndex(df, colName)
	if colIndex == -1 {
		return ""
	}

	return df.Elem(row, colIndex).String()
}

func getStringValueFromMap(data map[string]string, colName string, exists bool) string {
	if !exists {
		return ""
	}

	val, ok := data[colName]
	if !ok {
		return ""
	}

	return val
}

func getColumnIndex(df dataframe.DataFrame, colName string) int {
	for i, name := range df.Names() {
		if name == colName {
			return i
		}
	}
	return -1
}
