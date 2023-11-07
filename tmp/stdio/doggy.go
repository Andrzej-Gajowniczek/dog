package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var data string

	// Check if a filename is provided as a command-line argument
	if len(os.Args) > 1 {
		filename := os.Args[1]
		fileData, err := readFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
		} else {
			data = fileData
		}
	} else {
		// Check if data is available from a pipe
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			// Data is available from a pipe
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := scanner.Text()
				// Process the line as needed
				data += "Received from pipe: " + line + "\n"
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			}
		}
	}

	if data != "" {
		fmt.Println("Data received:")
		fmt.Println(data)
	} else {
		// No pipe data and no filename provided, show help information
		fmt.Println("Usage: " + os.Args[0] + " [filename]")
	}
}

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content, nil
}
