package LSB

import "image/color"

func NewImageIterator(image *LsbImage, startX int) *ImageIterator {
	return &ImageIterator{
		LsbImage: image,
		currentX: startX,
		currentY: -1,
	}
}

type ImageIterator struct {
	LsbImage *LsbImage
	currentX int
	currentY int
}

func (iterator *ImageIterator) Value() (color.NRGBA, int, int) {
	x, y := iterator.currentX, iterator.currentY
	r, g, b, _ := iterator.LsbImage.Image.At(x, y).RGBA()
	return color.NRGBA{uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), 255}, x, y
}

func (iterator *ImageIterator) Next() bool {
	iterator.currentY++

	if iterator.currentY < iterator.LsbImage.MaxY {
		return true
	}
	iterator.currentY = 0
	iterator.currentX++

	if iterator.currentX < iterator.LsbImage.MaxX {
		return true
	}

	return false
}
