package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/sys/unix"
)

type RGB struct {
	r uint8
	g uint8
	b uint8
	_ uint8
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
}

func main() {
	term.InitScreen()
	term.ClearScreen()
	term.RawMode()
	term.CursorHide()
	term.GetSize()
	term.CreateBlockBuffer()
	term.RestoreNormal()
	term.CursorShow()

	reader, err := os.Open("ocelot.jpg")
	defer reader.Close()
	if err != nil {
		log.Fatal("can't open img", err)
	}
	origImg, _, err := image.Decode(reader)
	origImg = origImg
	newImg := resize.Resize(uint(term.xBlock), uint(term.yBlock), origImg, resize.Lanczos3)

	bound := newImg.Bounds()
	ximg := bound.Max.X
	yimg := bound.Max.Y
	term.CursorAt(50, 10)
	fmt.Printf("ximg:%d yimg:%d\n", ximg, yimg)

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
