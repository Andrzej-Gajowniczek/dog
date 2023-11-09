package main

import (
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
	fd int
	//	mode            int
	xMax   int
	yMax   int
	xBlock int
	yBlock int
	//	raw             bool
	//	activeBuffer    uint8
	//	flushBuffer     uint8
	//	clearBuffer     uint8
	blockBuffer *[]RGB
	//	originalTermios unix.Termios
	Lflag           uint32
	termProportions float64
	//imgProportions  float64
	xImgResized int
	yImgResized int
	scanColors  *[]RGB
	man         map[RGB]int
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
		{0x17, 0x14, 0x21, 0},
		{0xbd, 0x1b, 0x21, 1},
		{0x26, 0xa2, 0x69, 2},
		{0xff, 0x74, 0x00, 3},
		{0x12, 0x48, 0x8b, 4},
		{0xa3, 0x47, 0xba, 5},
		{0x2a, 0xa1, 0xb3, 6},
		{0xd0, 0xcf, 0xcc, 7},
		{0x5e, 0x5c, 0x64, 8},
		{0xf6, 0x61, 0x51, 9},
		{0x33, 0xda, 0x7a, 10},
		{0xe9, 0xad, 0x0c, 11},
		{0x2a, 0x7b, 0xde, 12},
		{0xc0, 0x61, 0xcb, 13},
		{0x33, 0xc7, 0xde, 14},
		{0xff, 0xff, 0xff, 15},
	}
	term.scanColors = &list

	term.GetSize()
	term.InitScreen()

	term.ClearScreen()
	//term.RawMode()
	//term.CursorHide()

	//term.RestoreNormal()
	term.CursorShow()
	//	reader, err := os.Open("ocelot.jpg")
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
	//term.CursorAt(50, 10)
	//fmt.Printf("ximg:%d yimg:%d\n", ximg, yimg)

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
	term.RenderBlockGfxFrameGray()
	term.RenderBlockGfxFrame256()
	term.render8()
	term.RenderBlockGfxFrame8()
	term.RenderBlockGfxFrame808()
	shades()
	/*
		time.Sleep(4000 * time.Microsecond)
		term.RenderBlockGfxFrame256()
		time.Sleep(4000 * time.Microsecond)
	*/
	//}
	/*	time.Sleep(1 / 60 * time.Second)time.Sleep(1 * time.Second)
		term.RenderBlockGfxFrameGray()
		time.Sleep(1 * time.Second)
		term.RenderBlockGfxFrameRGB()
	*/ //time.Sleep(1 * time.Second)
	//term.CursorShow()
	//		fmt.Printf("type:%T", newImg)

	//		term.CursorAt(40, 20)
	//		fmt.Printf("len of buffer: %d\n", len(*term.blockBuffer))

	//fmt.Println("tproportions:", term.termProportions)
	//	fmt.Printf("xBlock:%d yBlock:%d\n", term.xBlock, term.yBlock)

	// fmt.Println(manual)
	//	term.manualMapping()
}
