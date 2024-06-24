package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// Usage: your_sqlite3.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		header := make([]byte, 100)

		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}

		var pageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}

		schemaBuf := make([]byte, pageSize)
		_, err = databaseFile.Read(schemaBuf)
		if err != nil {
			log.Fatal(err)
		}
		re, _ := regexp.Compile("CREATE TABLE")
		tables := re.FindAll(schemaBuf, -1)
		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", len(tables))
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
