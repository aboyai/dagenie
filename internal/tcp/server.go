package tcp

import (
	"bufio"
	"dagenie/internal/dagdb"
	"dagenie/internal/dql"
	"fmt"
	"net"
	"strings"
)

const (
	green = "\033[32m"
	red   = "\033[31m"
	reset = "\033[0m"
)

// DAGenie banner using dashes
const dagBanner = `
----------- DAGenie LightSpeed Server -----------
`

func StartTCPServer(db *dagdb.DAGDB, address string) error {
	// 🟢 Print banner in green
	fmt.Println(green + dagBanner + reset)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("❌ Failed to start TCP server: %v", err)
	}
	fmt.Printf("🚀 Dagenie TCP server running at %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("❌ Connection error: %v\n", err)
			continue
		}
		go handleConnection(conn, db)
	}
}

/*func handleConnection(conn net.Conn, db *dagdb.DAGDB) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		queryLine, err := reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte("❌ Error reading query\n"))
			return
		}

		queryLine = strings.TrimSpace(queryLine)
		queryLine = strings.TrimSuffix(queryLine, ";") // Trim semicolon
		if strings.ToLower(queryLine) == "exit" {
			conn.Write([]byte("👋 Bye!\n"))
			return
		}

		fmt.Printf("📨 Received query: %s\n", queryLine)

		// Unified DQL handling
		result, err := dql.ExecuteDQL(db, queryLine)
		fmt.Println(result)

		if err != nil {
			conn.Write([]byte(fmt.Sprintf("❌ %v\n", err)))
		} else {
			conn.Write([]byte(result + "\n"))
		}

		// Prompt next query
		conn.Write([]byte("📥 Ready for next query...\n"))
	}
}
*/

/*func handleConnection(conn net.Conn, db *dagdb.DAGDB) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		queryLine, err := reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte("❌ Error reading query\n"))
			return
		}

		queryLine = strings.TrimSpace(queryLine)
		queryLine = strings.TrimSuffix(queryLine, ";")
		if strings.ToLower(queryLine) == "exit" {
			conn.Write([]byte("👋 Bye!\n"))
			return
		}

		// Unified handler
		result, err := dql.ExecuteDQL(db, queryLine)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("❌ %v\n", err)))
		} else {
			conn.Write([]byte(result + "\n"))
		}
		conn.Write([]byte("📥 Ready for next query...\n"))
	}
}
*/

func handleConnection(conn net.Conn, globalDB *dagdb.DAGDB) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var clientDB *dagdb.DAGDB = globalDB // default DB

	for {
		queryLine, err := reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte("❌ Error reading query\n"))
			return
		}

		queryLine = strings.TrimSpace(queryLine)
		queryLine = strings.TrimSuffix(queryLine, ";")
		if strings.ToLower(queryLine) == "exit" {
			conn.Write([]byte("👋 Bye!\n"))
			return
		}
		fmt.Printf("📨 Received query: %s\n", queryLine)

		// Pass client-specific DB
		result, newDB, err := dql.ExecuteDQLWithContext(clientDB, queryLine)
		if newDB != nil {
			// Close old clientDB if needed
			if clientDB != globalDB {
				clientDB.Close()
			}
			clientDB = newDB
		}

		if err != nil {
			conn.Write([]byte(fmt.Sprintf("❌ %v\n", err)))
		} else {
			conn.Write([]byte(result + "\n"))

		}
		conn.Write([]byte("📥 Ready for next query...\n"))
	}
}
