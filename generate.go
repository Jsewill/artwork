package artwork

import (
	"image"
	"image/draw"
	"sort"
)

/*import (
	nftstorage "github.com/nftstorage/go-client"
	_ "github.com/stretchr/testify/assert"
	_ "golang.org/x/net/context"
	_ "golang.org/x/oauth2"
)

// Just something to keep go happy about the import until we decide whether or not to use it.
func Upload() {
	nft := nftstorage.NewNFT()
}*/

// @TODO: Rework the WeightMap functions, either make them more generic, or make them based on an attribute type that is fairly generic/useful.

// Sum returns the sum of the weight map.
func (a AttributeWeightMap) Sum() (sum float64) {
	for _, w := range a {
		sum += w
	}
	return
}

/* Computes an array of intervals which represent the normalized, weighted distribution of an AttributeMap. The sum of the receiver, "a" should equal 1.0.
Here's a visual:

[]AttributeWeightMap{aV:0.2, aW:0.15, aX:0.25, aY:0.1, aZ:0.3}

Would become something like,

[]AttributeWeightIntervalMap{aV:0.2, aW:0.35, aX:0.60, aY:0.7, aZ:1.0}

Order is irrelevant.

This function returns the distribution as an *AttributeWeightInterval slice
*/
func (a AttributeWeightMap) Intervals() AttributeWeightIntervals {
	// Make sure we have more than one weight.
	if len(a) <= 1 {
		// @TODO: Maybe this should return nil?
		//return nil
	}
	// Make the CDF slice
	prevSum := 0.0
	i, ints := 0, make(AttributeWeightIntervals, len(a))
	for attr, w := range a {
		ints[i] = &AttributeWeightInterval{
			Weight:    w + prevSum,
			Attribute: attr,
		}
		prevSum += w
		i++
	}
	// Sort the CDF slice
	sort.Slice(ints, func(i, j int) bool {
		return ints[i].Weight < ints[j].Weight
	},
	)
	return ints
}

type AttributeWeightInterval struct {
	Attribute
	Weight float64
}

type AttributeWeightIntervals []*AttributeWeightInterval

// Attribute returns the first attribute for which f is less than its computed distribution threshold. It assumes itself to be sorted.
func (a AttributeWeightIntervals) Attribute(f float64) Attribute {
	for _, awi := range a {
		// Since slice is sorted, this is all we should need.
		if f < awi.Weight {
			return awi.Attribute
		}
	}
	return 0
}

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
