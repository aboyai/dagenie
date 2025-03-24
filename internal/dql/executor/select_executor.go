package executor

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"dagenie/internal/dagdb"
	"dagenie/internal/dql/ast"

	"github.com/olekukonko/tablewriter"
)

// Regex to detect aggregate functions like SUM(duration)
var aggregateRegex = regexp.MustCompile(`(?i)(sum|avg|max|min|count)\s*\(\s*([a-zA-Z0-9_*]+)\s*\)`)

func ExecuteSelect(db *dagdb.DAGDB, selectAST *ast.SelectQueryAST) (string, error) {
	fmt.Println("S1")
	if selectAST.Table != "dag" {
		return "", fmt.Errorf("❌ Unsupported table: %s", selectAST.Table)
	}

	validFields := map[string]bool{
		"_id": true, "dagid": true, "id": true, "name": true,
		"status": true, "payload": true, "dependencies": true,
		"duration": true, "retries": true,
	}

	// Expand SELECT *
	fields := selectAST.Fields
	if len(fields) == 1 && fields[0] == "*" {
		fields = []string{"id", "name", "status", "payload", "dependencies", "dagid", "duration", "retries", "_id"}
		selectAST.Fields = fields
	}

	// Validate fields (skip aggregates)
	for _, field := range fields {
		fieldLower := strings.ToLower(field)
		if aggregateRegex.MatchString(fieldLower) {
			continue
		}
		if !validFields[fieldLower] {
			return "", fmt.Errorf("❌ Unknown field: %s", field)
		}
	}
	fmt.Println("START")

	// Load tasks
	var tasks []dagdb.DAGTask
	var err error
	if objectID, ok := selectAST.Conditions["_id"]; ok {
		tasks, err = db.QueryByObjectID(objectID)
	} else if dagID, ok := selectAST.Conditions["dagid"]; ok {
		tasks, err = db.ListTasksByDAG(dagID)
	} else {
		tasks, err = db.ListAllTasks()
	}
	if err != nil {
		return "", fmt.Errorf("❌ Task fetch error: %v", err)
	}

	fmt.Println("Here1")

	// Filter by WHERE
	var filtered []dagdb.DAGTask
	for _, task := range tasks {
		match := true
		for condField, condVal := range selectAST.Conditions {
			val := strings.ToLower(getField(task, condField))
			if val != strings.ToLower(condVal) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, task)
		}
	}

	fmt.Println("Here2")

	// COUNT(*)
	if selectAST.IsCount && len(selectAST.Aggregates) == 0 {
		return fmt.Sprintf("🔢 Count=%d\n\033[32m✅ Done\033[0m", len(filtered)), nil
	}

	// Handle Aggregates
	if len(selectAST.Aggregates) > 0 {
		if len(selectAST.GroupBy) > 0 {
			return executeGroupedAggregates(filtered, selectAST)
		}
		return executeGlobalAggregates(filtered, selectAST)
	}

	// ORDER BY (non-aggregates only)
	if len(selectAST.OrderBy) > 0 {
		sort.SliceStable(filtered, func(i, j int) bool {
			for _, ob := range selectAST.OrderBy {
				vi := getField(filtered[i], ob.Field)
				vj := getField(filtered[j], ob.Field)
				if vi == vj {
					continue
				}
				if ob.Desc {
					return vi > vj
				}
				return vi < vj
			}
			return false
		})
	}

	// LIMIT
	if selectAST.Limit > 0 && len(filtered) > selectAST.Limit {
		filtered = filtered[:selectAST.Limit]
	}
	fmt.Println("Here3")

	// No results
	if len(filtered) == 0 {
		return "❌ No results", nil
	}

	fmt.Println("FINAL")
	// Final output
	return formatSelectResults(filtered, fields), nil
}

func formatSelectResults(tasks []dagdb.DAGTask, fields []string) string {
	var sb strings.Builder

	// Prepare headers
	headers := []string{}
	for _, field := range fields {
		fieldUpper := strings.ToUpper(field)
		if fieldUpper == "_ID" {
			fieldUpper = "ObjectID"
		}
		headers = append(headers, fieldUpper)
	}

	// Setup tablewriter
	table := tablewriter.NewWriter(&sb)
	table.SetHeader(headers)

	// Color settings
	headerColors := make([]tablewriter.Colors, len(headers))
	colColors := make([]tablewriter.Colors, len(headers))
	for i := range headers {
		headerColors[i] = tablewriter.Colors{tablewriter.FgHiWhiteColor}
		colColors[i] = tablewriter.Colors{tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(headerColors...)
	table.SetColumnColor(colColors...)
	table.SetBorder(true)

	// Build rows from task data
	for _, task := range tasks {
		row := []string{}
		for _, field := range fields {
			val := getField(task, field)
			row = append(row, val)
		}
		table.Append(row)
	}

	table.Render()
	sb.WriteString("\033[32m✅ Done\033[0m\n")
	return sb.String()
}

func executeGlobalAggregates(tasks []dagdb.DAGTask, ast *ast.SelectQueryAST) (string, error) {
	if len(tasks) == 0 {
		return "❌ No data to aggregate\n✅ Done", nil
	}

	var sb strings.Builder
	for _, agg := range ast.Aggregates {
		vals := getNumericFieldValues(tasks, agg.Field)

		switch agg.Func {
		case "SUM":
			sum := 0.0
			for _, v := range vals {
				sum += v
			}
			sb.WriteString(fmt.Sprintf("SUM(%s)=%.2f ", agg.Field, sum))

		case "AVG":
			if len(vals) == 0 {
				sb.WriteString(fmt.Sprintf("AVG(%s)=0.00 ", agg.Field))
			} else {
				sum := 0.0
				for _, v := range vals {
					sum += v
				}
				avg := sum / float64(len(vals))
				sb.WriteString(fmt.Sprintf("AVG(%s)=%.2f ", agg.Field, avg))
			}

		case "MAX":
			if len(vals) == 0 {
				sb.WriteString(fmt.Sprintf("MAX(%s)=N/A ", agg.Field))
			} else {
				max := vals[0]
				for _, v := range vals[1:] {
					if v > max {
						max = v
					}
				}
				sb.WriteString(fmt.Sprintf("MAX(%s)=%.2f ", agg.Field, max))
			}

		case "MIN":
			if len(vals) == 0 {
				sb.WriteString(fmt.Sprintf("MIN(%s)=N/A ", agg.Field))
			} else {
				min := vals[0]
				for _, v := range vals[1:] {
					if v < min {
						min = v
					}
				}
				sb.WriteString(fmt.Sprintf("MIN(%s)=%.2f ", agg.Field, min))
			}

		case "COUNT":
			sb.WriteString(fmt.Sprintf("COUNT(%s)=%d ", agg.Field, len(tasks)))
		}
	}

	sb.WriteString("\n✅ Done")
	return sb.String(), nil
}

func executeGroupedAggregates(tasks []dagdb.DAGTask, ast *ast.SelectQueryAST) (string, error) {
	if len(tasks) == 0 {
		return "❌ No data to group\n\033[32m✅ Done\033[0m", nil
	}

	type groupKey struct {
		values []string
		keyStr string
	}

	groupMap := make(map[string][]dagdb.DAGTask)
	groupKeys := []groupKey{}

	for _, task := range tasks {
		vals := []string{}
		for _, field := range ast.GroupBy {
			fmt.Printf("2: %s", getField(task, field))
			vals = append(vals, getField(task, field))
		}
		keyStr := strings.Join(vals, "||")
		if _, exists := groupMap[keyStr]; !exists {
			groupKeys = append(groupKeys, groupKey{values: vals, keyStr: keyStr})
		}
		groupMap[keyStr] = append(groupMap[keyStr], task)
	}

	// Prepare clean headers
	headers := []string{}
	for _, field := range ast.GroupBy {
		headers = append(headers, strings.ToUpper(field))
	}
	for _, agg := range ast.Aggregates {
		headers = append(headers, fmt.Sprintf("%s(%s)", agg.Func, strings.ToUpper(agg.Field)))
	}

	// Prepare rows
	var rows [][]string
	for _, g := range groupKeys {
		groupTasks := groupMap[g.keyStr]
		row := make([]string, 0, len(g.values))
		row = append(row, g.values...)

		for _, agg := range ast.Aggregates {
			vals := getNumericFieldValues(groupTasks, agg.Field)
			switch agg.Func {
			case "SUM":
				sum := 0.0
				for _, v := range vals {
					sum += v
				}
				row = append(row, fmt.Sprintf("%.2f", sum))
			case "AVG":
				sum := 0.0
				for _, v := range vals {
					sum += v
				}
				avg := 0.0
				if len(vals) > 0 {
					avg = sum / float64(len(vals))
				}
				row = append(row, fmt.Sprintf("%.2f", avg))
			case "MAX":
				max := 0.0
				if len(vals) > 0 {
					max = vals[0]
					for _, v := range vals[1:] {
						if v > max {
							max = v
						}
					}
				}
				row = append(row, fmt.Sprintf("%.2f", max))
			case "MIN":
				min := 0.0
				if len(vals) > 0 {
					min = vals[0]
					for _, v := range vals[1:] {
						if v < min {
							min = v
						}
					}
				}
				row = append(row, fmt.Sprintf("%.2f", min))
			case "COUNT":
				row = append(row, fmt.Sprintf("%d", len(groupTasks)))
			}
		}
		rows = append(rows, row)
	}

	// ORDER BY Aggregate
	if len(ast.OrderByAgg) > 0 {
		sort.SliceStable(rows, func(i, j int) bool {
			for _, ob := range ast.OrderByAgg {
				targetHeader := fmt.Sprintf("%s(%s)", ob.Func, strings.ToUpper(ob.Field))

				colIndex := -1
				for idx, hdr := range headers {
					if strings.EqualFold(hdr, targetHeader) {
						colIndex = idx
						break
					}
				}

				if colIndex == -1 || colIndex >= len(rows[i]) || colIndex >= len(rows[j]) {
					continue
				}

				vi := rows[i][colIndex]
				vj := rows[j][colIndex]
				var viFloat, vjFloat float64
				_, err1 := fmt.Sscanf(vi, "%f", &viFloat)
				_, err2 := fmt.Sscanf(vj, "%f", &vjFloat)
				if err1 != nil || err2 != nil {
					continue
				}
				if viFloat == vjFloat {
					continue
				}
				if ob.Desc {
					return viFloat > vjFloat
				}
				return viFloat < vjFloat
			}
			return false
		})
	}

	// LIMIT
	if ast.Limit > 0 && len(rows) > ast.Limit {
		rows = rows[:ast.Limit]
	}

	// Render table with colored fields
	var sb strings.Builder
	table := tablewriter.NewWriter(&sb)
	table.SetHeader(headers)

	// Add color: white headers, green fields
	headerColors := make([]tablewriter.Colors, len(headers))
	colColors := make([]tablewriter.Colors, len(headers))
	for i := range headers {
		headerColors[i] = tablewriter.Colors{tablewriter.FgHiWhiteColor}
		colColors[i] = tablewriter.Colors{tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(headerColors...)
	table.SetColumnColor(colColors...)

	for _, row := range rows {
		table.Append(row)
	}
	table.SetBorder(true)
	table.Render()

	sb.WriteString("\033[32m✅ Done\033[0m\n")
	return sb.String(), nil
}

/*
	func executeGroupedAggregates(tasks []dagdb.DAGTask, ast *ast.SelectQueryAST) (string, error) {
		if len(tasks) == 0 {
			return "❌ No data to group\n✅ Done", nil
		}

		// Group tasks
		groupMap := make(map[string][]dagdb.DAGTask)
		var groupKeys []string

		for _, task := range tasks {
			keyParts := []string{}
			for _, field := range ast.GroupBy {
				keyParts = append(keyParts, getField(task, field))
			}
			key := strings.Join(keyParts, "_")
			groupMap[key] = append(groupMap[key], task)
			if !contains(groupKeys, key) {
				groupKeys = append(groupKeys, key)
			}
		}

		// Prepare table headers
		headers := append(ast.GroupBy, []string{}...) // Group fields
		for _, agg := range ast.Aggregates {
			headers = append(headers, fmt.Sprintf("%s(%s)", agg.Func, agg.Field))
		}

		// Prepare table data
		var rows [][]string
		for _, key := range groupKeys {
			groupTasks := groupMap[key]
			keyParts := strings.Split(key, "_")
			row := append([]string{}, keyParts...) // Group field values

			for _, agg := range ast.Aggregates {
				vals := getNumericFieldValues(groupTasks, agg.Field)
				switch agg.Func {
				case "SUM":
					sum := 0.0
					for _, v := range vals {
						sum += v
					}
					row = append(row, fmt.Sprintf("%.2f", sum))
				case "AVG":
					sum := 0.0
					for _, v := range vals {
						sum += v
					}
					avg := 0.0
					if len(vals) > 0 {
						avg = sum / float64(len(vals))
					}
					row = append(row, fmt.Sprintf("%.2f", avg))
				case "MAX":
					max := 0.0
					if len(vals) > 0 {
						max = vals[0]
						for _, v := range vals[1:] {
							if v > max {
								max = v
							}
						}
					}
					row = append(row, fmt.Sprintf("%.2f", max))
				case "MIN":
					min := 0.0
					if len(vals) > 0 {
						min = vals[0]
						for _, v := range vals[1:] {
							if v < min {
								min = v
							}
						}
					}
					row = append(row, fmt.Sprintf("%.2f", min))
				case "COUNT":
					row = append(row, fmt.Sprintf("%d", len(groupTasks)))
				}
			}
			rows = append(rows, row)
		}

		// Build table as string
		var sb strings.Builder
		table := tablewriter.NewWriter(&sb)
		table.SetHeader(headers)
		for _, row := range rows {
			table.Append(row)
		}
		table.SetBorder(true)
		table.Render()

		sb.WriteString("\033[32m✅ Done\033[0m\n")
		return sb.String(), nil
	}
*/
func getNumericFieldValues(tasks []dagdb.DAGTask, field string) []float64 {
	var values []float64
	for _, task := range tasks {
		switch strings.ToLower(field) {
		case "duration":
			values = append(values, float64(task.Duration))
		case "retries":
			values = append(values, float64(task.Retries))
		}
	}
	return values
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
