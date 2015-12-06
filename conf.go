package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func readConfig(fName string) (map[string]string, error) {
	f, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	outMap := make(map[string]string)
	splitRe := regexp.MustCompile(`\s+`)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		p := splitRe.Split(line, 2)
		if len(p) < 2 {
			continue
		}
		outMap[strings.ToLower(p[0])] = p[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return outMap, nil
}
