/**
 * 源码来自: https://github.com/J-Siu/go-png2ico/blob/master/go-png2ico.go
 */
package convert

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// ICO structire
type ICO struct {
	file string
	fh   *os.File
}

// PNG structure
type PNG struct {
	file   string // filename
	fh     *os.File
	height uint8
	width  uint8
	depth  uint16 // bit/pixel
	size   uint32
	offset uint32
	isPNG  bool
	buf    []byte
}

// Open : open PNG file
func (png *PNG) Open(file string) error {
	fmt.Println("PNG:Open:", file)

	var e error
	var n int

	png.file = file
	png.isPNG = false

	png.fh, e = os.Open(file)
	if e != nil {
		return e
	}

	/*
		25byte PNG header - BigEndian
		00:	89 50 4e 47 0d 0a 1a 0a // 8byte - magic number
		IHDR chunk
		08:	xx xx xx xx // 4byte - chunk length
		12:	49 48 44 52 // 4byte - chunk type(IHDR)
		16:	xx xx xx xx // 4byte - width
		20:	xx xx xx xx // 4byte - height
		24:	xx          // 1byte - bit depth (bit/pixel)
	*/
	headerLen := 25
	header := make([]byte, headerLen)
	n, e = png.fh.Read(header)
	if e != nil {
		return e
	}
	fmt.Println("PNG:Open:Header:", hex.EncodeToString(header), "(", n, ")")

	// 8byte header[0:8] - magic number
	magic := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	if bytes.Equal(magic[:], header[:8]) {
		fmt.Println("PNG:Open: Found PNG magic")
	} else {
		return errors.New("not png")
	}

	// 4byte header[8:12] - chunk length - skipped

	// 4byte header[12:16] - chunk type IHDR
	if bytes.Equal([]byte("IHDR"), header[12:16]) {
		fmt.Println("PNG:Open: Found IHDR chunk")
	} else {
		return errors.New("PNG no IHDR chunk")
	}

	// It is PNG
	png.isPNG = true

	// 4byte header[16:20] - width
	width := binary.BigEndian.Uint32(header[16:20])
	fmt.Println("PNG:Open:width:", width)

	// 4byte header[20:24] - height
	height := binary.BigEndian.Uint32(header[20:24])
	fmt.Println("PNG:Open:height:", height)

	if width <= 256 && height <= 256 {
		// ICO format use 0 for 256px or larger
		if width >= 256 {
			width = 0
		}
		if height >= 256 {
			height = 0
		}
	}

	png.width = uint8(width)
	png.height = uint8(height)

	// 1byte header[25] - color depth
	png.depth = uint16(uint8(header[24]))
	fmt.Println("PNG:Open:depth:", png.depth)

	stat, _ := os.Stat(file)
	png.size = uint32(stat.Size())
	fmt.Println("PNG:Open:size:", png.size)

	// Pass all check, create PNG struct
	fmt.Println("PNG:Open:png:", *png)

	return nil
}

// Read : read PNG file
func (png *PNG) Read() error {
	fmt.Println("PNG:Read:", png.file)

	var e error
	var n int
	var n64 int64

	n64, e = png.fh.Seek(0, 0)
	if e != nil {
		fmt.Println("读取文件失败", e)
		return e
	}
	fmt.Println("PNG:Read:Seek:", n64)

	png.buf = make([]byte, png.size)
	n, e = png.fh.Read(png.buf)
	fmt.Println("PNG:Read:byte:", n)

	return e
}

// Open : open ICO filehandle
func (ico *ICO) Open(file string) error {
	var e error
	fmt.Println("ICO:Open:", file)
	ico.fh, e = os.Create(file)
	return (e)
}

// Write : write ICO
func (ico *ICO) Write(b *[]byte) error {
	var e error
	var n int
	n, e = ico.fh.Write(*b)
	fmt.Println("ICO:Write:byte:", n)
	return (e)
}

// ICONDIR - return ICONDIR byte array
func (ico *ICO) ICONDIR(num uint16) *[]byte {
	/*
		6byte ICONDIR - LittleEndian
		00:   00 00 // 2byte, must be 0
		02:   01 00 // 2byte, 1 for ICO
		04:   xx xx // 2byte, img number
	*/

	b := []byte{0, 0, 1, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[4:6], num)
	fmt.Println("ICO:ICONDIR:", hex.EncodeToString(b))
	return &b
}

// ICONDIRENTRY - return ICONDIRENTRY byte array
func (png *PNG) ICONDIRENTRY() *[]byte {
	fmt.Println("PNG:ICONDIRENTRY:png:", *png)
	/*
		16byte ICONDIRENTRY - LittleEndian
		00:   xx    // 1byte, width
		01:   xx    // 1byte, height
		02:   00    // 1byte, color palette number, 0 for PNG
		03:   00    // 1byte, reserved, always 0
		04:   00 00 // 2byte, color planes, 0 for PNG
		06:   xx xx // 2byte, color depth
		08:   xx xx xx xx // 4byte, image size
		12:   xx xx xx xx // 4byte, image offset
	*/

	b := make([]byte, 16)

	copy(b[0:6], []byte{png.width, png.height, 0, 0, 0, 0})
	binary.LittleEndian.PutUint16(b[6:8], png.depth)
	binary.LittleEndian.PutUint32(b[8:12], png.size)
	binary.LittleEndian.PutUint32(b[12:16], png.offset)
	fmt.Println("PNG:ICONDIRENTRY:", hex.EncodeToString(b))

	return &b
}

// func main() {
// 	//Debug

// 	// ARGs
// 	args := os.Args[1:]
// 	argc := len(args)
// 	switch argc {
// 	case 0:
// 		usage()
// 		os.Exit(0)
// 	case 1:
// 		helper.ErrCheck(errors.New("Input/Output file missing"))
// 	}

// 	fileout := args[argc-1]

// 	// Make sure destination file is *not* PNG
// 	png := new(PNG)
// 	if png.Open(fileout) == nil || png.isPNG {
// 		helper.ErrCheck(errors.New("Output file (" + png.file + ") is a PNG file."))
// 	} else {
// 		fmt.Println("main:", png.file, "not PNG")
// 	}

// 	// Get and calculate all PNGs info
// 	pngs := []*PNG{}
// 	pngc := argc - 1
// 	var pngTotalSize uint32 = 0
// 	var LenICONDIR uint32 = 6
// 	var LenICONDIRENTRY uint32 = 16
// 	var LenAllICONDIRENTRY uint32 = LenICONDIRENTRY * uint32(pngc)
// 	for i := 0; i < pngc; i++ {
// 		png := new(PNG)
// 		helper.ErrCheck(png.Open(args[i]))
// 		// offset = len(ICONDIR) + len(all ICONDIRENTRY) + len(all PNG before current one)
// 		png.offset = LenICONDIR + LenAllICONDIRENTRY + pngTotalSize
// 		pngs = append(pngs, png)
// 		pngTotalSize += png.size
// 	}

// 	// Open ICON
// 	ico := new(ICO)
// 	helper.ErrCheck(ico.Open(fileout))
// 	helper.ErrCheck(ico.Write(ico.ICONDIR(uint16(pngc))))
// 	// Write ICONDIRENTRY
// 	for i := 0; i < pngc; i++ {
// 		helper.ErrCheck(ico.Write(pngs[i].ICONDIRENTRY()))
// 	}
// 	// Copy PNG
// 	for i := 0; i < pngc; i++ {
// 		helper.ErrCheck(pngs[i].Read())
// 		helper.ErrCheck(ico.Write(&pngs[i].buf))
// 	}
// }
