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

var (
	ColorPrinter = struct {
		Red      string
		Yellow   string
		Colorize func(string, string) string
	}{
		Red:    "\033[31m",
		Yellow: "\033[33m",
		Colorize: func(color string, text string) string {
			if runtime.GOOS == "windows" {
				return text
			}
			colorReset := "\033[0m"
			return color + text + colorReset
		},
	}
)

const (
	DummyUsage = "########"
)

var (
	// Define short parameters ( don't set default value ).
	paramsOutputPath              = flag.String("o", "", DummyUsage)
	paramsBeforeImagePath         = flag.String("b", "", DummyUsage)
	paramsAfterImagePath          = flag.String("a", "", DummyUsage)
	paramsErrorIfDifferenceExists = flag.Bool("e", false, "\n"+DummyUsage)
	paramsNotOutputIfSameImage    = flag.Bool("n", false, "\n"+DummyUsage)
	paramsHelp                    = flag.Bool("h", false, "\n"+DummyUsage)
)

func init() {
	// Define long parameters and description ( set default value here if you need ).
	// Required parameters
	flag.StringVar(paramsOutputPath /*         */, "output" /*                  */, "" /*    */, ColorPrinter.Colorize(ColorPrinter.Yellow, "[required]")+" Output path of difference image")
	flag.StringVar(paramsBeforeImagePath /*    */, "before-image-path" /*       */, "" /*    */, ColorPrinter.Colorize(ColorPrinter.Yellow, "[required]")+" Image path (before)")
	flag.StringVar(paramsAfterImagePath /*     */, "after-image-path" /*        */, "" /*    */, ColorPrinter.Colorize(ColorPrinter.Yellow, "[required]")+" Image path (after)")
	// Optional parameters
	flag.BoolVar(paramsErrorIfDifferenceExists /*  */, "error-if-difference" /*  */, false /* */, "Be regarded as an error (status code 1) if difference exists")
	flag.BoolVar(paramsNotOutputIfSameImage /*    */, "not-output-if-same" /*  */, false /* */, "Not output difference image if inputs are same")
	flag.BoolVar(paramsHelp /*                 */, "help" /*                    */, false /* */, "Show help")
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
	if *paramsHelp {
		flag.Usage()
		os.Exit(0)
	}
	if *paramsOutputPath == "" || *paramsBeforeImagePath == "" || *paramsAfterImagePath == "" {
		fmt.Println(ColorPrinter.Colorize(ColorPrinter.Red, "[ERROR]") + " Missing required parameter.")
		flag.Usage()
		os.Exit(1)
	}

	img1 := mustLoadImage(*paramsBeforeImagePath)
	img2 := mustLoadImage(*paramsAfterImagePath)

	diffImg, deletions, insertions, equals := diffimage.DiffImage(img1, img2)

	println(fmt.Sprintf("%-15s", "deletions"), deletions)
	println(fmt.Sprintf("%-15s", "insertions"), insertions)
	println(fmt.Sprintf("%-15s", "equals"), equals)

	hasDifference := deletions != 0 || insertions != 0

	// Not output difference image if inputs are same.
	if *paramsNotOutputIfSameImage && !hasDifference {
		os.Exit(0)
	}

	mustSaveImage(diffImg, *paramsOutputPath)

	// Exit with error status if difference of 2 images exists.
	if *paramsErrorIfDifferenceExists && hasDifference {
		os.Exit(1)
	}

	os.Exit(0)
}
