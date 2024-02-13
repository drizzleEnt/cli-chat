package token

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadRefresh() (string, error) {
	file, err := os.Open("bin/token.txt")
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from file: %v", err)
	}

	return line, nil
}
