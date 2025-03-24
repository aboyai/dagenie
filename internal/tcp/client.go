package tcp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("SELECT",
		readline.PcItem("*"),
		readline.PcItem("id"),
		readline.PcItem("name"),
		readline.PcItem("status"),
		readline.PcItem("payload"),
		readline.PcItem("dependencies"),
		readline.PcItem("dagid"),
		readline.PcItem("_id"),
	),
	readline.PcItem("INSERT"),
	readline.PcItem("UPDATE"),
	readline.PcItem("DELETE"),
	readline.PcItem("EXIT"),
)

type dynamicCompleter struct{}

func StartTCPClient(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("🔗 Connected to %s\n", serverAddr)

	currentDB := "default"

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            fmt.Sprintf("📝 [%s] DQL > ", currentDB),
		HistoryFile:       "/tmp/dagenie_history.tmp",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		AutoComplete:      dynamicCompleter{},
		HistorySearchFold: true,
	})
	if err != nil {
		fmt.Printf("❌ Readline init error: %v\n", err)
		return
	}
	defer rl.Close()

	serverReader := bufio.NewReader(conn)
	var buffer strings.Builder

	for {
		rl.SetPrompt(fmt.Sprintf("📝 [%s] DQL > ", currentDB))
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("👋 Exiting.")
			break
		}
		line = strings.TrimSpace(line)

		if strings.ToLower(line) == "exit" {
			conn.Write([]byte("exit\n"))
			break
		}

		buffer.WriteString(line + " ")

		if strings.HasSuffix(line, ";") {
			query := strings.TrimSuffix(buffer.String(), ";")
			buffer.Reset()

			_, err = conn.Write([]byte(query + "\n"))
			if err != nil {
				fmt.Printf("❌ Write error: %v\n", err)
				break
			}

			var responseLines []string
			for {
				resp, err := serverReader.ReadString('\n')
				if err != nil {
					fmt.Printf("❌ Read error: %v\n", err)
					break
				}
				resp = strings.TrimRight(resp, "\r\n")

				// Update DB prompt if DB switched
				if strings.HasPrefix(resp, "✅ Using database") {
					parts := strings.Split(resp, "'")
					if len(parts) >= 2 {
						currentDB = parts[1]
						rl.SetPrompt(fmt.Sprintf("📝 [%s] DQL > ", currentDB))
					}
				}

				// Wait until we get 📥 Ready for next query... before proceeding
				if resp == "📥 Ready for next query..." {
					break // Stop reading; ready for next user query
				}

				// Collect response for display
				responseLines = append(responseLines, resp)
			}

			// Detect key=value lines for table display
			hasKeyValue := false
			for _, line := range responseLines {
				if strings.Contains(line, "=") && !strings.HasPrefix(line, "✅") && !strings.HasPrefix(line, "❌") {
					hasKeyValue = true
					break
				}
			}

			// Print result
			if hasKeyValue {
				printAsTable(responseLines)
			} else {
				for _, line := range responseLines {
					if strings.HasPrefix(line, "✅") {
						fmt.Println("\033[32m" + line + "\033[0m")
					} else if strings.HasPrefix(line, "❌") {
						fmt.Println("\033[31m" + line + "\033[0m")
					} else {
						fmt.Println(line)
					}
				}
			}

		}
	}
}

func printAsTable(lines []string) {
	var dataLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.ContainsAny(line, "📥✅❌📂👋") {
			continue
		}
		dataLines = append(dataLines, line)
	}

	if len(dataLines) == 0 {
		return
	}

	// Extract headers from the first line
	fieldOrder := []string{}
	headers := []string{}
	firstFields := strings.Fields(dataLines[0])
	for _, pair := range firstFields {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		fieldOrder = append(fieldOrder, key)
		if strings.ToLower(key) == "_id" {
			headers = append(headers, "ObjectID")
		} else {
			headers = append(headers, strings.ToUpper(key))
		}
	}

	// Table setup
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	// Colors
	headerColors := make([]tablewriter.Colors, len(headers))
	colColors := make([]tablewriter.Colors, len(headers))
	for i := range headers {
		headerColors[i] = tablewriter.Colors{tablewriter.FgHiWhiteColor}
		colColors[i] = tablewriter.Colors{tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(headerColors...)
	table.SetColumnColor(colColors...)
	table.SetBorder(true)

	// Append rows
	for _, line := range dataLines {
		fields := strings.Fields(line)
		fieldMap := make(map[string]string)
		for _, pair := range fields {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				fieldMap[kv[0]] = kv[1]
			}
		}
		row := []string{}
		for _, key := range fieldOrder {
			row = append(row, fieldMap[key])
		}
		table.Append(row)
	}

	table.Render()
}

func (c dynamicCompleter) Do(line []rune, pos int) ([][]rune, int) {
	input := strings.ToLower(string(line[:pos]))
	suggestions := []string{}

	if strings.HasPrefix(input, "select") || strings.Contains(input, "select") {
		suggestions = []string{"id", "name", "status", "payload", "dependencies", "dagid", "_id", "duration", "retries"}
	} else if strings.HasPrefix(input, "insert") {
		suggestions = []string{"insert into dag (...)", "values (...)"}
	} else if strings.HasPrefix(input, "update") {
		suggestions = []string{"update dag set", "where"}
	} else if strings.HasPrefix(input, "delete") {
		suggestions = []string{"delete from dag where"}
	} else if strings.HasPrefix(input, "use") {
		suggestions = []string{"use default;", "use logs;", "use users;"}
	} else if strings.HasPrefix(input, "show") {
		suggestions = []string{"show databases;"}
	} else {
		suggestions = []string{"select", "insert", "update", "delete", "use", "show", "exit"}
	}

	completions := [][]rune{}
	for _, s := range suggestions {
		completions = append(completions, []rune(s))
	}

	return completions, 0
}

/*func printAsTable(lines []string) {
	var dataLines []string

	// Filter valid data lines
	for _, line := range lines {
		line = strings.TrimSpace(line)
		fmt.Println(line)
		if line == "" || strings.Contains(line, "📥") || strings.Contains(line, "✅") ||
			strings.Contains(line, "❌") || strings.Contains(line, "📂") || strings.Contains(line, "👋") {
			continue
		}
		dataLines = append(dataLines, line)
	}

	if len(dataLines) == 0 {
		return
	}

	// Parse headers from the first valid line
	headers := []string{}
	firstFields := strings.Fields(dataLines[0])
	for _, pair := range firstFields {
		kv := strings.SplitN(pair, "=", 2)
		fieldName := kv[0]

		// 🔥 Fix: extract only valid field names (handle ORDER BY artifacts)
		fieldNameClean := sanitizeFieldName(fieldName)

		if strings.ToLower(fieldNameClean) == "_id" {
			fieldNameClean = "ObjectID"
		}
		fmt.Println(fieldName)
		headers = append(headers, fieldNameClean)
	}

	// Initialize table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	// Colors
	headerColors := make([]tablewriter.Colors, len(headers))
	colColors := make([]tablewriter.Colors, len(headers))
	for i := range headers {
		headerColors[i] = tablewriter.Colors{tablewriter.FgHiWhiteColor}
		colColors[i] = tablewriter.Colors{tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(headerColors...)
	table.SetColumnColor(colColors...)
	table.SetBorder(true)

	// Append rows
	for _, line := range dataLines {
		fields := strings.Fields(line)
		values := []string{}
		for _, pair := range fields {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				values = append(values, kv[1])
			} else {
				values = append(values, "")
			}
		}
		table.Append(values)
	}

	table.Render()
}
*/
