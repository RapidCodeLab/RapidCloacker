package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) (map[string][]IPItem, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make(map[string][]IPItem)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		splittedLine := strings.FieldsFunc(line, Split)

		item := IPItem{}

		itemIsIP := net.ParseIP(line)

		if itemIsIP != nil {
			item.IP = itemIsIP
			lines[splittedLine[0]+splittedLine[1]+splittedLine[2]] = append(lines[splittedLine[0]+splittedLine[1]+splittedLine[2]], item)
			continue
		}

		_, itemNet, err := net.ParseCIDR(line)
		if err != nil {
			log.Printf("[ERROR] %+v", err)
		}
		item.IPNet = itemNet

		lines[splittedLine[0]+splittedLine[1]+splittedLine[2]] = append(lines[splittedLine[0]+splittedLine[1]+splittedLine[2]], item)
	}

	return lines, scanner.Err()
}
