package main

import (
	"fmt"
	"os"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("nie podano parametrów")
		os.Exit(0)

	} else {
		fmt.Println("podano:", args[1:])

		for _, path := range args[1:] {
			_, err := os.Stat(path)
			if err != nil {
				fmt.Println(path, " isn't a file")
			} else {
				fmt.Println(path, "is a file")
			}

		}
	}
}
