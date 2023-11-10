package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

type RGB struct {
	r uint8
	g uint8
	b uint8
	n uint8
}
type terminal struct {
	fd              int
	xMax            int
	yMax            int
	xBlock          int
	yBlock          int
	blockBuffer     *[]RGB
	Lflag           uint32
	termProportions float64
	xImgResized     int
	yImgResized     int
	scanColors      *[]RGB
	man             map[RGB]int
	colorMatrix     map[RGB]string
}

func main() {

	manual := make(map[RGB]int, 27)

	manual[RGB{0, 0, 0, 0}] = 0
	manual[RGB{1, 0, 0, 0}] = 1
	manual[RGB{2, 0, 0, 0}] = 9
	manual[RGB{0, 1, 0, 0}] = 2
	manual[RGB{0, 2, 0, 0}] = 10
	manual[RGB{0, 0, 1, 0}] = 4
	manual[RGB{0, 0, 2, 0}] = 12
	manual[RGB{1, 1, 0, 0}] = 3
	manual[RGB{2, 1, 0, 0}] = 3
	manual[RGB{1, 2, 0, 0}] = 11
	manual[RGB{2, 2, 0, 0}] = 11
	manual[RGB{1, 0, 1, 0}] = 5
	manual[RGB{2, 0, 1, 0}] = 13
	manual[RGB{1, 0, 2, 0}] = 13
	manual[RGB{2, 0, 2, 0}] = 13
	manual[RGB{0, 1, 1, 0}] = 6
	manual[RGB{0, 2, 1, 0}] = 14
	manual[RGB{0, 1, 2, 0}] = 14
	manual[RGB{0, 2, 2, 0}] = 14
	manual[RGB{1, 1, 1, 0}] = 8
	manual[RGB{2, 1, 1, 0}] = 9
	manual[RGB{1, 2, 1, 0}] = 10
	manual[RGB{1, 1, 2, 0}] = 12
	manual[RGB{2, 2, 1, 0}] = 7
	manual[RGB{2, 1, 2, 0}] = 13
	manual[RGB{1, 2, 2, 0}] = 14
	manual[RGB{2, 2, 2, 0}] = 15
	term.man = manual
	//color2index := make(map[RGB]int, 16)
	list := []RGB{
		//	{0x17, 0x14, 0x21, 0},
		{0x10, 0x0c, 0x1c, 0},
		{0xde, 0x00, 0x13, 1},
		{0x00, 0xa7, 0x63, 2},
		{0xb2, 0x71, 0x42, 3},
		{0x00, 0x47, 0x72, 4},
		{0xb7, 0x39, 0xc2, 5},
		{0x00, 0xa5, 0xb8, 6},
		{0xd2, 0xd0, 0xcd, 7},
		{0x5e, 0x5b, 0x65, 8},
		{0xff, 0x51, 0x44, 9},
		{0x00, 0xe0, 0x6e, 10},
		{0xfd, 0xac, 0x00, 11},
		{0x00, 0x7d, 0xe7, 12},
		{0xd6, 0x57, 0xd3, 13},
		{0x00, 0xcc, 0xe3, 14},
		{0xff, 0xff, 0xff, 15},
	}
	term.scanColors = &list
	percentage := []float64{100, 50, 21.875}
	//matrix map[RGB]sting
	matrix := make(map[RGB]string, 768)
	term.colorMatrix = matrix
	for fg := 0; fg < 16; fg++ {

		fgRGB := list[fg]
		fgR := float64(fgRGB.r)
		fgG := float64(fgRGB.g)
		fgB := float64(fgRGB.b)
		for bg := 0; bg < 16; bg++ {
			bgRGB := list[bg]
			bgR := float64(bgRGB.r)
			bgG := float64(bgRGB.g)
			bgB := float64(bgRGB.b)

			for inx, ch := range []rune{'█', '▒', '░'} {
				//for inx, ch := range []rune{'▒'} {

				eR := fgR*percentage[inx]/100 + bgR*(100-percentage[inx])/100
				eG := fgG*percentage[inx]/100 + bgG*(100-percentage[inx])/100
				eB := fgB*percentage[inx]/100 + bgB*(100-percentage[inx])/100

				eRGB := RGB{uint8(eR), uint8(eG), uint8(eB), 0}

				matrix[eRGB] = fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%c", fg, bg, ch)

			}
		}
	}

	term.GetSize()
	term.InitScreen()

	term.ClearScreen()
	//term.RawMode()
	//term.CursorHide()

	//term.RestoreNormal()
	term.CursorShow()
	//reader, err := os.Open("ocelot.jpg")
	reader, err := os.Open("dog2.png")

	if err != nil {
		log.Fatal("can't open img", err)
	}
	defer reader.Close()
	origImg, _, _ := image.Decode(reader)
	//origImg = origImg
	newImg := resize.Resize(uint(term.xBlock), uint(0), origImg, resize.Lanczos3)

	bound := newImg.Bounds()
	ximg := bound.Max.X
	term.xImgResized = bound.Max.X
	yimg := bound.Max.Y
	term.yImgResized = bound.Max.Y

	term.CreateBlockBuffer()

	term.CursorAt(0, 0)
	for i := 0; i < ximg; i++ {
		for j := 0; j < yimg; j++ {
			r, g, b, _ := newImg.At(i, j).RGBA()
			(*term.blockBuffer)[j*term.xBlock+i].r = uint8(r >> 8)
			(*term.blockBuffer)[j*term.xBlock+i].g = uint8(g >> 8)
			(*term.blockBuffer)[j*term.xBlock+i].b = uint8(b >> 8)

		}
	}
	//for {
	term.RenderBlockGfxFrameRGB()
	//term.RenderBlockGfxFrameGray()
	//t
	term.renderMagic16()
	term.RenderBlockGfxFrame256()
	term.RenderBlockGfxFrameGray()
	term.renderBlueMoon()
	//term.RenderBlockGfxFrame8()
	//term.RenderBlockGfxFrame808()
	//shades()

	//fmt.Println("len(matrix)", len(matrix))
	/*
		var lastRGBerr float64 = 1000
		var bestRGB RGB
		lastRGB := RGB{255, 0, 0, 0}
		bestRGB = RGB{0, 0, 0, 0}

		var sign string = ""

		for {

			for rgb, ch := range matrix {
				lastR := float64(lastRGB.r)
				lastG := float64(lastRGB.g)
				lastB := float64(lastRGB.b)
				currR := float64(rgb.r)
				currG := float64(rgb.g)
				currB := float64(rgb.b)
				errRGB := math.Sqrt((lastR-currR)*(lastR-currR) + (lastG-currG)*(lastG-currG) + (lastB-currB)*(lastB-currB))
				if lastRGBerr > errRGB {
					lastRGBerr = errRGB
					bestRGB.r = rgb.r
					bestRGB.g = rgb.g
					bestRGB.b = rgb.b

					sign = ch

				}
			}
			lastRGB = bestRGB

			fmt.Printf("%s%s%s%s%s%s%s%s\x1b[m%03d,%03d,%03d-% 8f\n", sign, sign, sign, sign, sign, sign, sign, sign, bestRGB.r, bestRGB.g, bestRGB.b, lastRGBerr)
			lastRGBerr = 1000
			delete(matrix, RGB{uint8(bestRGB.r), uint8(bestRGB.g), uint8(bestRGB.b), 0})
			if len(matrix) == 0 {
				break
			}
		}
	*/
	//	term.renderMagic()

}
