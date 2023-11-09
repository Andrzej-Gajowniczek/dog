package main

import (
	"fmt"
	"os"
)

type mixol struct {
	fg int
	bg int
	ch rune
}

func matrix() *[]string {

	shadding := []rune{'█', '▒', '░'}
	//grays := []uint8{0, 0, 4, 12, 6, 14, 14, 7, 15}
	//limitG := len(grays)
	/*
		plik, err := os.OpenFile("color-info.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)

		if err != nil {
			log.Fatal("nie można otworzyć pliku")
		}
		defer plik.Close()
	*/
	var koloriada []string
	index := 0
	for i := 0; i < 16; i++ {
		//bcont := true
		for j := 0; j < 16; j++ {
			if i == j {
				continue
			}
			for _, ch := range shadding {

				p := fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", i, j, ch)
				index++
				koloriada = append(koloriada, p)
				os.Stdout.WriteString(p)

			}
			//println("\x1b[0m\n")
			//fmt.Fprintf(plik, "%s,", p)
			//plik.WriteString(",")

			/*	if !bcont {
				break
			}*/
		}
		fmt.Println("\x1b[0m")
	}
	println("\x1b[0m\n")
	return &koloriada
}

func main() {

	a := matrix()
	//a = a

	for _, str := range *a {
		os.Stdout.WriteString(str)
	}
	println("\x1b[0m\n")
}
