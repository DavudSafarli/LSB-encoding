package LSB

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type LsbImage struct {
	path  string
	Image *image.RGBA

	// metadata is binary representation of length of secret
	// metadataBitsLength is number of bits required to represent length of text
	MetadataBitsLength int
	MaxTextLength      int

	MaxX int
	MaxY int
}

// NewLsbImage creates a new instance of LsbImage
func NewLsbImage(path string) *LsbImage {
	reader, _ := os.Open(path)
	ext := filepath.Ext(path)
	var readonlyImage image.Image
	var err error
	if ext != ".png" {
		panic("cannot embed secret to jpg images, since they get compressed")
	}
	readonlyImage, err = png.Decode(reader)
	check(err)

	bounds := readonlyImage.Bounds()

	pixelCount := bounds.Max.X * bounds.Max.Y
	metadataBitsLength := len(strconv.FormatInt(int64(pixelCount), 2))
	maxTextLength := (pixelCount - metadataBitsLength) / 8
	mutableImage := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	draw.Draw(mutableImage, bounds, readonlyImage, image.Point{0, 0}, draw.Src)
	return &LsbImage{
		path:               path,
		MaxX:               bounds.Max.X,
		MaxY:               bounds.Max.Y,
		MetadataBitsLength: metadataBitsLength,
		MaxTextLength:      maxTextLength,
		Image:              mutableImage,
	}
}

func (lsbImage *LsbImage) Save() {
	dir, nameWithExt := filepath.Split(lsbImage.path)
	ext := filepath.Ext(lsbImage.path)
	filename := nameWithExt[0 : len(nameWithExt)-len(ext)]
	newImagePath := dir + filename + "_embedded" + ext
	fg, err := os.Create(newImagePath)
	check(err)
	defer fg.Close()
	if ext != ".png" {
		panic("cannot embed secret to jpg images, since they get compressed")
	}
	err = png.Encode(fg, lsbImage.Image)
	check(err)
}

// CreateEmbeddableImage creates a new image with all Blue values EVEN
func CreateEmbeddableImage(sourceImage string) *LsbImage {
	lsbImage := NewLsbImage(sourceImage)

	for x := 0; x < lsbImage.MaxX; x++ {
		for y := 0; y < lsbImage.MaxY; y++ {
			pixel := lsbImage.Image.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			b := originalColor.B
			if b%2 == 1 {
				b--
			}
			c := color.RGBA{
				R: originalColor.R,
				G: originalColor.G,
				B: b,
				A: originalColor.A,
			}
			lsbImage.Image.SetRGBA(x, y, c)
		}
	}

	return lsbImage
}
