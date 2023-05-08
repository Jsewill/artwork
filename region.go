package artwork

import "image"

// Region defines an overlay region, with center-point coordinates, Coords, a slice of applicable kinds, Kinds, and an asset to overlay, Asset.
type Region struct {
	*Asset
	Coords *image.Point
	Kinds  []string
	Scale  *Scale
	//Transform f64.Aff3
}

func NewRegion() *Region {
	return &Region{
		Kinds: make([]string, 0),
	}
}

// Coordinates is a getter function for region coordinates. Defaults to center.
func (r *Region) Coordinates() *image.Point {
	// Default to center.
	if r.Coords == nil {
		r.Coords = &image.Point{r.Bounds().Min.Add(r.Bounds().Dx() / 2), r.Bounds().Min.Add(r.Bounds().Dy() / 2)}
	}
	return r.Coords
}
