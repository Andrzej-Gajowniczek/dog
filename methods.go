package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"golang.org/x/sys/unix"
)

var term terminal

func (t *terminal) CursorHide() {
	os.Stdout.WriteString("\x1b[?25l")
}
func (t *terminal) CursorShow() {
	os.Stdout.WriteString("\x1b[?25h\x1b[m")
}
func (t *terminal) ClearScreen() {
	//fmt.Print("\033[H\033[2J")
	os.Stdout.WriteString("\x1b[H\x1b[2J")
}

func (t *terminal) CursorAt(x, y int) {
	positionCode := fmt.Sprintf("\033[%d;%dH", y, x)
	os.Stdout.WriteString(positionCode)

}

func (t *terminal) Print(s string) {

}
func (t *terminal) PrintAt216(x, y int, rf, gf, bf, rb, gb, bb uint8, s string) {

}
func (t *terminal) PrintAtRGB(x, y int, rf, gf, bf, rb, gb, bb uint8, s string) {
	p := fmt.Sprintf("\x1b[%d,%dH\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s", y, x, rf, gf, bf, rb, gb, bb, s)
	os.Stdout.WriteString(p)
}
func (t *terminal) printAtGray(x, y int, i, j uint8, s string) {

}

func (t *terminal) InitScreen() {
	t.fd = int(os.Stdout.Fd())
}

func (t *terminal) RawMode() {

	termios, err := unix.IoctlGetTermios(term.fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}
	t.Lflag = termios.Lflag
	// Zmień ustawienia terminala na "raw mode"
	termios.Lflag &^= unix.ICANON | unix.ECHO
	err = unix.IoctlSetTermios(term.fd, unix.TCSETS, termios)
	if err != nil {
		panic(err)
	}
}

func (t *terminal) RestoreNormal() {
	termios, err := unix.IoctlGetTermios(term.fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}
	termios.Lflag = t.Lflag
	err = unix.IoctlSetTermios(term.fd, unix.TCSETS, termios)
	if err != nil {
		panic(err)
	}

}

func (t *terminal) GetSize() {

	ws, err := unix.IoctlGetWinsize(t.fd, unix.TIOCGWINSZ)
	if err != nil {
		log.Fatalln("I can't get size of terminal", err)
		t.xMax = 0
		t.yMax = 0

	}
	t.xMax = int(ws.Col)
	t.yMax = int(ws.Row)
	t.xBlock = t.xMax
	t.yBlock = t.yMax << 1
	t.termProportions = float64(t.xBlock) / float64(t.yBlock)

}

func (t *terminal) CreateBlockBuffer() {
	blocks := t.xBlock * t.yBlock

	t.blockBuffer = new([]RGB)
	*t.blockBuffer = make([]RGB, blocks)
}

func (t *terminal) RenderBlockGfxFrameRGB() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xBlock
	ySize := t.yBlock
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]
			p = fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%c", v.r, v.g, v.b, w.r, w.g, w.b, feedBlock)
			os.Stdout.WriteString(p)
		}
	}

}

func (t *terminal) RenderBlockGfxFrame256() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xBlock
	ySize := t.yBlock
	divider := uint8(47)
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]

			v.r = addRandomNoise256(v.r)
			v.g = addRandomNoise256(v.g)
			v.b = addRandomNoise256(v.b)

			w.r = addRandomNoise256(w.r)
			w.g = addRandomNoise256(w.g)
			w.b = addRandomNoise256(w.b)

			//v, w = t.applySteinbergHoffmanDithering(v, w, float32(divider), x, y)

			c := 16 + 36*(v.r/divider) + 6*(v.g/divider) + 1*(v.b/divider)
			k := 16 + 36*(w.r/divider) + 6*(w.g/divider) + 1*(w.b/divider)
			p = fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", c, k, feedBlock)
			os.Stdout.WriteString(p)

		}
	}
}

func (t *terminal) RenderBlockGfxFrameGray() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xBlock
	ySize := t.yBlock
	var divider float32 = 11
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]

			var vv float32 = (0.299*float32(v.r) + 0.587*float32(v.g) + 0.114*float32(v.g)) / divider
			var ww float32 = (0.299*float32(w.r) + 0.587*float32(w.g) + 0.114*float32(w.g)) / divider
			vv = addRandomNoise(vv)
			ww = addRandomNoise(ww)
			c := 255 - (23 - uint8(vv))
			k := 255 - (23 - uint8(ww))
			p = fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", c, k, feedBlock)
			os.Stdout.WriteString(p)
		}
	}
}

func addRandomNoise(value float32) float32 {
	// Define the noise level (adjust as needed)
	noiseLevel := .5

	// Generate random noise in the range [-noiseLevel/2, noiseLevel/2]
	noise := (rand.Float32() - 0.5) * float32(noiseLevel)

	// Add the noise to the original value
	return value + noise
}

func addRandomNoise256(value uint8) uint8 {
	// Define the noise level (adjust as needed)
	noiseLevel := uint8(11)

	// Generate random noise in the range [-noiseLevel/2, noiseLevel/2]
	noise := rand.Intn(int(noiseLevel+1)) - int(noiseLevel/2)

	// Add the noise to the original value
	result := int(value) + noise

	// Ensure the result stays within the valid range (0-255)
	if result < 0 {
		return 0
	} else if result > 255 {
		return 255
	}

	return uint8(result)
}

/*
func (t *terminal) applySteinbergHoffmanDithering(v RGB, w RGB, divider float32, x, y int) (RGB, RGB) {
	// Define the error coefficients for Steinberg-Hoffman dithering
	errCoefs := [3][3]float32{
		{0.0, 0.0, 7.0 / 16.0},
		{3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0},
		{0.0, 1.0 / 16.0, 0.0},
	}

	// Calculate the errors for v and w
	var errorV RGB
	var errorW RGB
	errorV.r = v.r - uint8((float32(v.r)/255.0)*divider*255.0)
	errorV.g = v.g - uint8((float32(v.g)/255.0)*divider*255.0)
	errorV.b = v.b - uint8((float32(v.b)/255.0)*divider*255.0)

	errorW.r = w.r - uint8((float32(w.r)/255.0)*divider*255.0)
	errorW.g = w.g - uint8((float32(w.g)/255.0)*divider*255.0)
	errorW.b = w.b - uint8((float32(w.b)/255.0)*divider*255.0)

	// Apply dithering error diffusion to v
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if x+i < t.xBlock && y+j < t.yBlock {
				(*t.blockBuffer)[(y+j)*t.xBlock+(x+i)].r += uint8(errCoefs[i][j] * float32(errorV.r))
				(*t.blockBuffer)[(y+j)*t.xBlock+(x+i)].g += uint8(errCoefs[i][j] * float32(errorV.g))
				(*t.blockBuffer)[(y+j)*t.xBlock+(x+i)].b += uint8(errCoefs[i][j] * float32(errorV.b))
			}
		}
	}

	// Apply dithering error diffusion to w (on the next row)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if x+i < t.xBlock && y+j+1 < t.yBlock {
				(*t.blockBuffer)[(y+j+1)*t.xBlock+(x+i)].r += uint8(errCoefs[i][j] * float32(errorW.r))
				(*t.blockBuffer)[(y+j+1)*t.xBlock+(x+i)].g += uint8(errCoefs[i][j] * float32(errorW.g))
				(*t.blockBuffer)[(y+j+1)*t.xBlock+(x+i)].b += uint8(errCoefs[i][j] * float32(errorW.b))
			}
		}
	}

	return v, w
}
*/

func (t *terminal) render8() {
	/*	kolorGray := []string{"\x1b[38;5;0m\x1b[48;5;8m█", "\x1b[38;5;0m\x1b[48;5;8m▒", "\x1b[38;5;0m\x1b[48;5;8m░", "\x1b[38;5;8m\x1b[48;5;7m█",
		"\x1b[38;5;8m\x1b[48;5;7m▒", "\x1b[38;5;8m\x1b[48;5;7m░", "\x1b[38;5;7m\x1b[48;5;15m█", "\x1b[38;5;7m\x1b[48;5;15m▒",
		"\x1b[38;5;7m\x1b[48;5;15m░", "\x1b[38;5;15m\x1b[48;5;15m█", "\x1b[38;5;15m\x1b[48;5;15m▒", "\x1b[38;5;15m\x1b[48;5;15m░"}
	*/ //for _, c := range kolorGray {
	//	fmt.Print(c)
	//}

	//	fmt.Println("\x1b[m", len(kolorGray))

	//t.CursorAt(0, 0)
	//feedBlock := '\u2580'

	kolorGray := *shades()
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xBlock
	ySize := t.yBlock
	var divider float32 = 11
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]

			var vv float32 = (0.299*float32(v.r) + 0.587*float32(v.g) + 0.114*float32(v.g)) / divider
			var ww float32 = (0.299*float32(w.r) + 0.587*float32(w.g) + 0.114*float32(w.g)) / divider
			vv = addRandomNoise(vv)
			ww = addRandomNoise(ww)
			//vv = vv
			//ww = ww
			k := ww*0.3333333333 + vv*0.66666666666
			p = kolorGray[int(k)+1]
			os.Stdout.WriteString(p)
		}
	}

	//fmt.Printf("\x1b[m;\nx:%d*y:%d\n", x, y)

}

func shades() *[]string {

	shadding := []rune{'█', '▒', '░'}
	grays := []uint8{0, 0, 4, 12, 6, 14, 14, 7, 15}
	//grays := []int8{0, 0, 1, 5, 13, 9, 11, 7, 15}
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