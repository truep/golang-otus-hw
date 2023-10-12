package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	ev := make(Environment)
	for _, file := range files {
		fi, err := file.Info()
		if err != nil {
			continue
		}

		if !strings.Contains(file.Name(), "=") || !file.IsDir() {
			if fi.Size() == 0 {
				ev[file.Name()] = EnvValue{
					Value:      "",
					NeedRemove: true,
				}
				continue
			}

			filePath := fmt.Sprintf("%s/%s", dir, file.Name())
			fileContent, err := readFirstLine(filePath)
			if err != nil {
				continue
			}

			value := normalizeValue(fileContent)

			ev[file.Name()] = EnvValue{
				Value:      value,
				NeedRemove: false,
			}
		}
	}

	return ev, nil
}

func readFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", nil
}

func normalizeValue(value string) string {
	value = strings.TrimRight(value, "\t ")
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value
}
