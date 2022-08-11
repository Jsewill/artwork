package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/draw"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
)

// Asset is a type of image asset which has a type, Kind, a filepath, Path, and a slice of available overlay regions, Regions.
type Asset struct {
	Kind    string
	Path    string
	Image   image.Image // @TODO: Consider embedding.
	Parent  *Region
	Regions []*Region
}

// Load loads an Asset's image into *Asset.Image. Returns an error if something went wrong along the way.
func (a *Asset) Load() error {
	// Check for a path.
	if a.Path == "" {
		err := fmt.Errorf("Failed to load asset image: Path is empty.")
		logErr.Println(err)
		return err
	}
	// Open the image file.
	ifile, err := os.Open(a.Path)
	if err != nil {
		err := fmt.Errorf("Failed to load asset image: %s", err)
		logErr.Println(err)
		return err
	}
	// Decode the image file. @TODO: consider allowing for selecting a frame from a GIF.
	a.Image, format, err = image.Decode(ifile)
	if err != nil {
		err = fmt.Errorf("Failed to load asset image: Error while decoding %q: %s", a.Path, err)
		logErr.Println(err)
		return err
	}

	return nil
}

// IsLoaded reports whether an *Asset.Image is not nil. Returns true if not nil, otherwise returns false.
func (a *Asset) IsLoaded() bool {
	if a.Image == nil {
		return false
	}
	return true
}

// Region defines an overlay region, with center-point coordinates, Coords, a slice of applicable kinds, Kinds, and an asset to overlay, Asset.
type Region struct {
	*Asset
	Coords *image.Point
	Kinds  []string
	Scale  *Scale
	//Transform f64.Aff3
}

// Coordinates is a getter function for region coordinates. Defaults to center.
func (r *Region) Coordinates() *image.Point {
	// Default to center.
	if r.Coords == nil {
		r.Coords = &image.Point{r.Bounds().Min.Add(r.Bounds().Dx() / 2), r.Bounds().Min.Add(r.Bounds().Dy() / 2)}
	}
	return r.Coords
}

// Composite climbs the current composition tree branch and composites down from the leaves. Returns image.Image, nil, when successful; nil, nil, at leaf; and nil/image.Image, error when there was a failure. @TODO: Find things to error about...
func (a *Asset) Composite() (image.Image, error) {
	// Check asset. If nil, this region is a leaf.
	if a == nil || a.Regions == nil {
		// Nothing to composite.
		return nil, nil
	}
	// Not a leaf
	// Get the current asset's bounds, well-formed.
	abounds := a.Image.Bounds().Canon()
	// Scale the branch asset if scale not zero or negative. @TODO: Consider implementing negative scaling for image inversion. Better yet, implement affine transformation.
	if a.Parent.Scale != nil {
		abounds := ScaleRectangle(a.Parent.Scale, orig)
	}
	// Create a new canvas to draw on and pass down the tree. We redraw the canvas to preserve the asset image. Fortunately, this happens outside the region loop.
	canvas := image.NewNRGBA(abounds)
	// Draw, and potentially scale, this branch asset onto the canvas.
	draw.ApproxBiLinear.Scale(canvas, canvas.Bounds(), a.Image, a.Image.Bounds(), draw.Over, nil)
	// Climb the tree.
	for _, region := range a.Regions {
		comp, err := region.Composite()
		// Composition failed farther along the branch. Send it on down the tree.
		if err != nil {
			/*err := fmt.Errorf("Error compositing asset: %+v: %s", a, err)
			logErr.Println(err)*/
			return comp, err
		}
		// Was at leaf. Send this image back down the tree.
		if comp == nil {
			return canvas, nil
		}
		// Composite this asset with the branch composite.
		// Get the branch composite bounds, well-formed.
		cbounds := comp.Bounds().Canon()
		// Expand the current canvas if necessary.
		if !cbounds.In(canvas.Bounds()) {
			// Replace the old canvas with the new one.
			canvas = GrowImage(canvas, cbounds)
		}
		// Center branch composite image on region coordinates.
		offset := CenterOffset(region.Coordinates(), cbounds)
		// Draw the asset onto the canvas.
		draw.Over.Draw(canvas, canvas.Bounds(), comp, cbounds.Min.Add(offset))
	}

	return canvas, nil
}
