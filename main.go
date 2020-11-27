package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lokks307/go-emf/emf"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

const VERSION = "0.0.1"

func main() {

	// Flag

	logDebugFlag := flag.Bool("debug", false, "print out debug message")
	inFile := flag.String("in", "", "emf file to convert")
	outFile := flag.String("out", "./out.png", "png file to output")

	flag.Parse()

	// Logger

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "15:04:05.000",
		FullTimestamp:   true,
		ForceColors:     true,
	})

	if *logDebugFlag {
		log.SetLevel(log.TraceLevel)
	}

	log.SetOutput(colorable.NewColorableStdout())

	// file

	if *inFile == "" {
		fmt.Println("")
		fmt.Println("GO-EMF: EMF images converter (ver. ", VERSION, ")")
		fmt.Println("")
		fmt.Println("Usage: ./go-emf [options]")
		flag.PrintDefaults()
		fmt.Println("")
		os.Exit(0)
		return
	}

	cwd, err := os.Getwd()

	if err != nil {
		os.Exit(0)
		return
	}

	inFilePath := filepath.Join(cwd, *inFile)

	var fdata []byte

	fdata, err = ioutil.ReadFile(inFilePath)
	if err != nil {
		log.Errorf("no such file %s", inFilePath)
		os.Exit(0)
		return
	}

	log.Info("EMF file reading...")
	emfFile := emf.ReadFile(fdata)
	log.Info("EMF file reading... done")

	log.Info("Converting EMF file to PNG...")
	emfFile.DrawToPNG(*outFile)
	log.Info("Converting EMF file to PNG... done")

}
