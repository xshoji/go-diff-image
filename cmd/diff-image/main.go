package main

import (
	"bytes"
	"flag"
	"fmt"
	diffimage "github.com/xshoji/go-diff-image"
	"image"
	"image/png"
	"log"
	"os"
	"regexp"
	"runtime"
)

func mustOpen(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func mustLoadImage(filename string) image.Image {
	f := mustOpen(filename)
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func mustSaveImage(img image.Image, output string) {
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(f, img)
}

func rate16(c uint32, r float64) uint16 {
	fc := float64(c)
	return uint16(65535 - (65535-fc)*r)
}

const (
	DummyUsage = "########"
)

var (
	TerminalColorReset  = "\033[0m"
	TerminalColorYellow = "\033[33m"
	// Define short parameters ( don't set default value ).
	paramsOutputPath      = flag.String("o", "", DummyUsage)
	paramsBeforeImagePath = flag.String("b", "", DummyUsage)
	paramsAfterImagePath  = flag.String("a", "", DummyUsage)
	paramsHelp            = flag.Bool("h", false, "\n"+DummyUsage)
)

func init() {

	if runtime.GOOS == "windows" {
		TerminalColorReset = ""
		TerminalColorYellow = ""
	}

	// Define long parameters and description ( set default value here if you need ).
	// Required parameters
	flag.StringVar(paramsOutputPath /*      */, "output" /*            */, "" /*     */, TerminalColorYellow+"[required]"+TerminalColorReset+" Output path")
	flag.StringVar(paramsBeforeImagePath /* */, "before-image-path" /* */, "" /*     */, TerminalColorYellow+"[required]"+TerminalColorReset+" Image path (before)")
	flag.StringVar(paramsAfterImagePath /*  */, "after-image-path" /*   */, "" /*     */, TerminalColorYellow+"[required]"+TerminalColorReset+" Image path (after)")
	// Optional parameters
	flag.BoolVar(paramsHelp /*              */, "help" /*              */, false /*  */, "Show help")
}

func main() {

	// Set DummyUsage
	b := new(bytes.Buffer)
	flag.CommandLine.SetOutput(b)
	flag.Usage()
	re := regexp.MustCompile("(-\\S+)( *\\S*)+\n*\\s+" + DummyUsage + "\n*\\s+(-\\S+)( *\\S*)+\n")
	usage := re.ReplaceAllString(b.String(), "  $1, -$3$4\n")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), usage)
	}

	flag.Parse()
	if *paramsHelp || *paramsOutputPath == "" || *paramsBeforeImagePath == "" || *paramsAfterImagePath == "" {
		flag.Usage()
		os.Exit(0)
	}

	img1 := mustLoadImage(*paramsBeforeImagePath)
	img2 := mustLoadImage(*paramsAfterImagePath)

	dst := diffimage.DiffImage(img1, img2)

	mustSaveImage(dst, *paramsOutputPath)
}
