package main

import (
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

func shades() *[]string {

	shadding := []rune{'█', '▒', '░'}
	grays := []uint8{0, 0, 4, 12, 6, 14, 14, 7, 15}
	limitG := len(grays)
	plik, err := os.OpenFile("color-info.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)

	if err != nil {
		log.Fatal("nie można otworzyć pliku")
	}
	defer plik.Close()

	var koloriada []string
	for m, i := range grays {
		//bcont := true
		for _, j := range shadding {
			k := m + 1
			if k == limitG {
				k = m
				//	bcont = false
			}
			c := grays[k]

			p := fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", i, c, j)
			koloriada = append(koloriada, p)
			//fmt.Fprintf(plik, "%s,", p)
			//plik.WriteString(",")
			//os.Stdout.WriteString(p)
			/*	if !bcont {
				break
			}*/
		}
	}
	return &koloriada
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: programname <unicode_character>")
		return
	}

	inputChar := os.Args[1]

	// Convert the input character to its Unicode code point
	runeValue, _ := utf8.DecodeRuneInString(inputChar)

	// Format the Unicode code point as \uNumber
	unicodeCode := fmt.Sprintf("\\u%04X", runeValue)

	fmt.Printf("Character: %s\nUnicode Code: %s\n", inputChar, unicodeCode)

	rainbow := shades()

	for i, v := range *rainbow {
		fmt.Printf("kolor:%s%d\x1b[m\n", v, i)
	}
}
