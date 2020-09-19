package main

import (
	"flag"
	"fmt"
	"image/color"
	_ "image/png"
	"os"
	"strconv"

	"github.com/DavudSafarli/LSB-encoding/LSB"
)

func main() {
	// Subcommands
	embedCommand := flag.NewFlagSet("embed", flag.ExitOnError)
	extractCommand := flag.NewFlagSet("extract", flag.ExitOnError)

	// embed subcommand flag pointers
	embedTextPtr := embedCommand.String("text", "", "Text to encode. (Required)")
	embedPathPtr := embedCommand.String("path", "", "PNG Image to embed text to. (Required)")
	extractPathPtr := extractCommand.String("path", "", "Path to the image. (Required)")

	if len(os.Args) < 2 {
		fmt.Errorf("extract or embed sub-command is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "extract":
		extractCommand.Parse(os.Args[2:])
	case "embed":
		embedCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if extractCommand.Parsed() {
		if *extractPathPtr == "" {
			extractCommand.PrintDefaults()
			os.Exit(1)
		}
		secret := extract(*extractPathPtr)
		fmt.Println(secret)
	}

	if embedCommand.Parsed() {
		if *embedTextPtr == "" || *embedPathPtr == "" {
			embedCommand.PrintDefaults()
			os.Exit(1)
		}
		embed(*embedTextPtr, *embedPathPtr)
	}
}

func embed(text string, sourceImage string) {
	secret := LSB.Secret{Text: text}
	lsbImage := LSB.CreateEmbeddableImage(sourceImage)

	metadataAsBitSlice := secret.GetMetadataAsBitSlice(lsbImage.MetadataBitsLength)
	textAsASCIIBits := secret.TextAsASCIIBits()

	bitsToWrite := append(metadataAsBitSlice, textAsASCIIBits...)

	pixelIterator := LSB.NewImageIterator(lsbImage, 0)
	for _, b := range bitsToWrite {
		pixelIterator.Next()
		if !b {
			continue
		}
		pixel, x, y := pixelIterator.Value()

		if pixel.B%2 == 1 {
			panic("source image shouldn't had a pixel with Blue being odd")
		}
		lsbImage.Image.SetRGBA(x, y, color.RGBA{
			R: pixel.R,
			G: pixel.G,
			B: pixel.B + 1,
			A: pixel.A,
		})
	}

	lsbImage.Save()
}

func extract(secretHiddenImagePath string) string {
	lsbImage := LSB.NewLsbImage(secretHiddenImagePath)

	imgIterator := LSB.NewImageIterator(lsbImage, 0)

	metadata := ""
	for i := 0; i < lsbImage.MetadataBitsLength; i++ {
		imgIterator.Next()

		pixel, _, _ := imgIterator.Value()

		if pixel.B%2 == 1 {
			metadata += "1"
			continue
		}
		metadata += "0"
	}
	secretLength, _ := strconv.ParseInt(metadata, 2, 64)

	secretTextAsciiBits := make([]bool, 0, secretLength*8)
	for i := 0; i < int(secretLength)*8; i++ {
		imgIterator.Next()

		pixel, _, _ := imgIterator.Value()

		if pixel.B%2 == 1 {
			secretTextAsciiBits = append(secretTextAsciiBits, true)
			continue
		}
		secretTextAsciiBits = append(secretTextAsciiBits, false)
	}
	secretText := LSB.ConvertAsciiBitsToText(secretTextAsciiBits)

	return secretText
}
