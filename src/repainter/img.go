package repainter

import (
	"image"
	"image/color"
	"image/draw"
)

type pixel struct {
	x      int
	y      int
	curClr color.NRGBA
}

type Img struct {
	repaintByte uint8
	rgbaBlack   color.NRGBA
	rgbaWhite   color.NRGBA
}

func NewImage() *Img {
	return &Img{
		repaintByte: 210, //affects on recolour accuracy. Closer to 255 - the numbers in the result image will be thinner
		rgbaBlack:   color.NRGBA{R: 0, G: 0, B: 0, A: 0},
		rgbaWhite:   color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func (r *Img) Repaint(img image.Image, inverted bool) image.Image {
	drawImg := img.(draw.Image)

	bounds := drawImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			clr := img.At(x, y).(color.NRGBA)

			if inverted {
				r.recolourPixelInverted(drawImg, pixel{x: x, y: y, curClr: clr})
			} else {
				r.recolourPixel(drawImg, pixel{x: x, y: y, curClr: clr})
			}
		}
	}

	return drawImg
}

func (r *Img) recolourPixel(img draw.Image, p pixel) {
	if r.isPixelLookLikeWhite(p.curClr) {
		img.Set(p.x, p.y, r.rgbaWhite)
		return
	}

	img.Set(p.x, p.y, r.rgbaBlack)
}

func (r *Img) recolourPixelInverted(img draw.Image, p pixel) {
	if r.isPixelLookLikeWhite(p.curClr) {
		img.Set(p.x, p.y, r.rgbaBlack)
		return
	}

	img.Set(p.x, p.y, r.rgbaWhite)
}

func (r *Img) isPixelLookLikeWhite(clr color.NRGBA) bool {
	return clr.R > r.repaintByte && clr.G > r.repaintByte && clr.B > r.repaintByte && clr.A > r.repaintByte
}
