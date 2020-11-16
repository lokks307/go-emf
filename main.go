package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lokks307/emftoimg/emf"
)

const VERSION = "0.1.0"

var errlog = log.New(os.Stderr, "emf: ", 0)

var (
	flagVersion = flag.Bool("version", false, "")
)

var usage = `EMF images converter

Usage: emftoimg [inputfile]
   	--version  print the version number

`

func main() {

	var fdata []byte
	var err error

	fdata, err = ioutil.ReadFile("convert.emf")

	t1 := time.Now()
	file, err := emf.ReadFile(fdata)
	if err != nil {
		errlog.Fatal(err)
	}
	e1 := time.Since(t1)

	t2 := time.Now()
	img := file.Draw()
	e2 := time.Since(t2)

	var f io.Writer

	f, err = os.Create("out" + ".png")
	if err != nil {
		errlog.Fatal(err)
	}
	defer f.(*os.File).Close()

	draw2dimg.SaveToPngFile("out.png", img)

	// err = png.Encode(f, img)
	// if err != nil {
	// 	errlog.Fatal(err)
	// }

	errlog.Printf("file %d bytes reading %.3f ms conversion %.3f ms\n",
		len(fdata),
		float64(e1.Nanoseconds())/1000000,
		float64(e2.Nanoseconds())/1000000)

}

func isatty(fd uintptr) bool {
	return true
	// var termios syscall.Termios

	// _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd,
	// 	uintptr(TCGETS),
	// 	uintptr(unsafe.Pointer(&termios)),
	// 	0,
	// 	0,
	// 	0)
	// return err == 0
}
