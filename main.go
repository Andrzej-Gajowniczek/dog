package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/sys/unix"
)

type RGB struct {
	r uint8
	g uint8
	b uint8
	n uint8
}
type terminal struct {
	fd              int
	mode            int
	xMax            int
	yMax            int
	xBlock          int
	yBlock          int
	raw             bool
	activeBuffer    uint8
	flushBuffer     uint8
	clearBuffer     uint8
	blockBuffer     *[]RGB
	originalTermios unix.Termios
	Lflag           uint32
	termProportions float64
	imgProportions  float64
	xImgResized     int
	yImgResized     int
	scanColors      *[]RGB
}

func main() {

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
	reader, err := os.Open("ocelot.jpg")
	//reader, err := os.Open("dog2.png")
	defer reader.Close()
	if err != nil {
		log.Fatal("can't open img", err)
	}
	origImg, _, err := image.Decode(reader)
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

}
