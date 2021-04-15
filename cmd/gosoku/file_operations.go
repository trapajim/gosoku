package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func FindLine(s string) int {
	f, err := os.Open("./main.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), s) {
			return line
		}

		line++
	}
	return 1
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}
