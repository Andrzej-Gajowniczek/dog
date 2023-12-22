package main

import (
	"fmt"
	"image"
	_ "image/gif"
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
	xImg            int
	yImg            int
	imgRatio        float64
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
	term.CursorAt(0, 0)
	term.ClearScreen()
	term.CursorShow()

	//reader, err := os.Open("ocelot.jpg")
	reader, err := os.Open("dog2.png")
	if err != nil {
		log.Fatal("can't open img", err)
	}
	defer reader.Close()
	origImg, _, _ := image.Decode(reader)
	term.imgRatio = CountImgRatio(origImg)
	term.xImg, term.yImg = GetImgSize(origImg)
	virtualY := GetYsizeBasedOnXandRatio(term.xBlock, term.imgRatio)
	newImg := resize.Resize(uint(term.xBlock), uint(virtualY), origImg, resize.Lanczos3)
	//term.CreateBlockBuffer()
	term.RenderBlockBuffer(newImg)
	term.RenderBlockGfxFrameRGB()
	term.RenderMagic16()
	term.RenderBlockGfxFrame256()
	term.RenderBlockGfxFrameGray()
	term.RenderBlueMoon()
}
