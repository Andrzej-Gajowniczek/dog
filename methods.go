package main

import (
	"fmt"
	"log"
	"math"
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

	rest := t.yImgResized % 2
	blocks := t.xImgResized * (t.yImgResized + rest)

	t.blockBuffer = new([]RGB)
	*t.blockBuffer = make([]RGB, blocks)
}

func (t *terminal) RenderBlockGfxFrameRGB() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xImgResized
	ySize := t.yImgResized
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]
			p = fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%c", v.r, v.g, v.b, w.r, w.g, w.b, feedBlock)
			os.Stdout.WriteString(p)
		}
		os.Stdout.WriteString("\x1b[m\n")
	}

}

func (t *terminal) RenderBlockGfxFrame256() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xImgResized
	ySize := t.yImgResized
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
	xSize := t.xImgResized
	ySize := t.yImgResized
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
		os.Stdout.WriteString("\x1b[m\n")
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

func (t *terminal) renderBlueMoon() {
	kolorGray := *shades()
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xImgResized
	ySize := t.yImgResized
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
		os.Stdout.WriteString("\x1b[m\n")
	}

}

func (t *terminal) renderMagic16() {
	//kolorGray := *shades()
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xImgResized
	ySize := t.yImgResized
	//var divider float32 = 11
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]
			//r := (v.r>>1 + w.r>>1)
			//g := (v.g>>1 + w.g>>1)
			//b := (v.b>>1 + w.b>>1)
			r := float64(v.r)
			g := float64(v.g)
			b := float64(v.g)
			r2 := float64(w.r)
			g2 := float64(w.g)
			b2 := float64(w.g)

			qr := (66.333333*r + 33.66666*r2) / 100
			qg := (66.333333*g + 33.66666*g2) / 100
			qb := (66.333333*b + 33.66666*b2) / 100
			challenger := RGB{uint8(qr), uint8(qg), uint8(qb), 0}

			//var vv float32 = (0.299*float32(v.r) + 0.587*float32(v.g) + 0.114*float32(v.g)) / divider
			//var ww float32 = (0.299*float32(w.r) + 0.587*float32(w.g) + 0.114*float32(w.g)) / divider
			//vv = addRandomNoise(vv)
			//ww = addRandomNoise(ww)
			//vv = vv
			//ww = ww
			//k := ww*0.3333333333 + vv*0.66666666666
			p = t.findMagicRGB(challenger)
			os.Stdout.WriteString(p)
		}
		os.Stdout.WriteString("\x1b[m\n")
	}
}

func (t *terminal) findMagicRGB(rgb RGB) string {
	var result string
	lastRGBerr := float64(1000)
	var bestRGB RGB

	for rgbn, ch := range t.colorMatrix {
		lastR := float64(rgb.r)
		lastG := float64(rgb.g)
		lastB := float64(rgb.b)
		currR := float64(rgbn.r)
		currG := float64(rgbn.g)
		currB := float64(rgbn.b)
		errRGB := math.Sqrt((lastR-currR)*(lastR-currR) + (lastG-currG)*(lastG-currG) + (lastB-currB)*(lastB-currB))
		if lastRGBerr > errRGB {
			lastRGBerr = errRGB
			bestRGB.r = rgb.r
			bestRGB.g = rgb.g
			bestRGB.b = rgb.b

			result = ch

		}
	}
	return result
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
		for _, j := range shadding {
			k := m + 1
			if k == limitG {
				k = m
			}
			c := grays[k]

			p := fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", i, c, j)
			koloriada = append(koloriada, p)

			//os.Stdout.WriteString(p)
		}

	}
	//	fmt.Println("\x1b[m")
	return &koloriada
}

func (t *terminal) leastSquares(rgb RGB) int {
	var color int
	var lastValue float64 = 500
	var newValue float64
	rr := float64(rgb.r)
	gg := float64(rgb.g)
	bb := float64(rgb.b)
	var irr, igg, ibb float64
	for i, iRGB := range *t.scanColors {
		irr = float64(iRGB.r)
		igg = float64(iRGB.g)
		ibb = float64(iRGB.b)
		newValue = math.Sqrt((rr-irr)*(rr-irr) + (gg-igg)*(gg-igg) + (bb-ibb)*(bb-ibb))
		if lastValue > newValue {
			lastValue = newValue
			color = i
		}
	}
	return color
}

func (t *terminal) RenderBlockGfxFrame8() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var bg, fg int
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xImgResized
	ySize := t.yImgResized
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]
			fg = t.leastSquares(v)
			bg = t.leastSquares(w)
			p = fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", fg, bg, feedBlock)
			os.Stdout.WriteString(p)
		}
		os.Stdout.WriteString("\x1b[m\n")
	}

}

func (t *terminal) RenderBlockGfxFrame808() {
	//t.CursorAt(0, 0)
	feedBlock := '\u2580'
	var bg, fg int
	var y, x int
	var v, w RGB
	var p string
	xSize := t.xImgResized
	ySize := t.yImgResized
	for y = 0; y < ySize; y += 2 {
		for x = 0; x < xSize; x++ {
			v = (*t.blockBuffer)[y*xSize+x]
			w = (*t.blockBuffer)[(y+1)*xSize+x]
			w.n = 0
			v.n = 0
			if v.r <= 17 {
				v.r = 0
			} else if v.r <= 23 {
				v.r = 1
			} else {
				v.r = 2
			}
			if v.g <= 17 {
				v.g = 0
			} else if v.g <= 47 {
				v.g = 1
			} else {
				v.g = 2
			}
			if v.b <= 17 {
				v.b = 0
			} else if v.b <= 63 {
				v.b = 1
			} else {
				v.b = 2
			}
			if w.r <= 17 {
				w.r = 0
			} else if w.r <= 32 {
				w.r = 1
			} else {
				w.r = 2
			}
			if w.g <= 17 {
				w.g = 0
			} else if w.g <= 78 {
				w.g = 1
			} else {
				w.g = 2
			}
			if v.b <= 17 {
				w.b = 0
			} else if v.b <= 90 {
				w.b = 1
			} else {
				w.b = 2
			}
			fg = t.man[v]
			bg = t.man[w]
			p = fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", fg, bg, feedBlock)
			os.Stdout.WriteString(p)
		}
		os.Stdout.WriteString("\x1b[m\n")
	}

}
