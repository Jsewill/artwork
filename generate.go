package artwork

import (
	"image"
	"image/draw"
)

/* @TODO: Find a replacement for NFTStorage. Make an interface for image storage service APIs.

import (
	nftstorage "github.com/nftstorage/go-client"
	_ "github.com/stretchr/testify/assert"
	_ "golang.org/x/net/context"
	_ "golang.org/x/oauth2"
)

// Just something to keep go happy about the import until we decide whether or not to use it.
func Upload() {
	nft := nftstorage.NewNFT()
}*/

// GrowImage enlarges an image to match the supplied rectangle bounds. Returns the grown image.Image.
func GrowImage(orig image.Image, bounds image.Rectangle) image.Image {
	// Only grow image if needed.
	if !cbounds.In(canvas.Bounds()) {
		// Get the union of our branch composite and the current canvas.
		union := bounds.Union(orig.Bounds())
		// Create a canvas to replace the old one.
		grown := image.NewNRGBA(union)
		// Draw the asset onto the canvas.
		draw.Over.Draw(grown, union, orig, union.Min)
		return grown
	}
	return orig
}

type Scale struct {
	X, Y float64
}

// ScaleRectangle scales a rectangle, orig, base on a scale factor, s. returns an image.Rectangle.
func ScaleRectangle(sx, sy float64, orig image.Rectangle) image.Rectangle {
	// Only scale if s is non-zero, non-negative. @TODO: Consider implementing negative scaling for image inversion.
	if sx == 0 && sy == 0 {
		return orig
	}
	// Normalize the scale factor. A scale factor of less than 1 should be subtracted from the max.
	var xFactor, yFactor float64
	if sx != 0 {
		xFactor = sx - 1.0
	}
	if sy != 0 {
		yFactor = sy - 1.0
	}
	// Create a point to add to the bounds max.
	pscale := image.Point{int(float64(orig.Dx()) * xFactor), int(float64(orig.Dy()) * yFactor)}
	// Translate orig to Zero.
	var omin image.Point
	if orig.Min != image.Point {
		omin := orig.Min
		orig = orig.Sub(omin)
	}
	// Add the scale factor to shrink or grow the rectangle.
	orig := image.Rectangle{image.Point{}, orig.Max.Add(pscale)}
	// Put orig back where it was.
	if omin != image.Point {
		orig = orig.Add(omin)
	}
	return orig
}

// CenterOffset gets the offset required to center a rectangle, bounds, on a point, center
func CenterOffset(center image.Point, bounds image.Rectangle) image.Point {
	return center.Sub(bounds.Max.Sub(bounds.Min).Div(2))
}
