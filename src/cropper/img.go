package cropper

import (
	"github.com/disintegration/imaging"
	"github.com/oliamb/cutter"
	"image"
)

func NewImage() *Img {
	return &Img{}
}

type Img struct {}

func (p *Img) CropAnswer(img image.Image, num int) (image.Image, error) {
	var y int

	switch num {
	case 1:
		y = 131
	case 2:
		y = 186
	case 3:
		y = 241
	}

	img, err := cutter.Crop(img, cutter.Config{
		Width:  50,
		Height: 25,
		Anchor: image.Point{X: 95, Y: y},
		Mode:   cutter.TopLeft, // optional, default value
	})
	if err != nil {
		return nil, err
	}

	img = imaging.AdjustGamma(img, 0.2)
	img = imaging.Sharpen(img, 4)
	img = imaging.AdjustContrast(img, 30)
	// img = imaging.Invert(img)

	return img, nil
}

func (p *Img) CropQuestion(img image.Image) (image.Image, error) {
	img, err := cutter.Crop(img, cutter.Config{
		Width:  40,
		Height: 20,
		Anchor: image.Point{X: 85, Y: 71},
		Mode:   cutter.TopLeft, // optional, default value
	})
	if err != nil {
		return nil, err
	}

	img = imaging.AdjustContrast(img, 70)
	img = imaging.Sharpen(img, 1)
	img = imaging.Invert(img)
	img = imaging.Grayscale(img)

	return img, nil
}

func (p *Img) CropCaptcha(img image.Image) (image.Image, error) {
	return cutter.Crop(img, cutter.Config{
		Width:  250,
		Height: 500,
		Mode:   cutter.Centered, // optional, default value
	})
}