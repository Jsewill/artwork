package artwork

import (
	"fmt"
	"image"
	"image/draw"
	"log"
)

// Piece represents a piece of artwork, with a *Region slice, Regions, for defining the composition tree, and an image.Image onto which it is to be composited, Canvas.
type Piece struct {
	Id uint
	// @TODO: add attributes and a way to set them. Perhaps, a func type.
	*Asset
}

// NewPiece creates a new piece from a base image.Image, canvas, or creates new image from bounds.
func NewPiece(id uint, canvas image.Image, bounds *image.Rectangle) *Piece {
	// Create canvas if it wasn't supplied.
	if canvas == nil {
		if bounds == nil {
			// Create canvas with zero size.
			canvas = image.NewNRGBA(image.Rectangle{})
		} else {
			// At least we got bounds. Create canvas from bounds.
			canvas = image.NewNRGBA(&bounds)
		}
	}
	return &Piece{Id: id, Asset: &Asset{Image: canvas}}
}

// Build creates an asset tree from a set of asset configuration data. Returns nil on succes, error on failure.
func (p *Piece) Build() error {
	// Build tree from Configuration.
	//
	return nil
}

// Composite walks Regions, attempting to composite the entire composition tree onto the canvas.
func (p *Piece) Composite() error {
	// Check for canvas.
	if p.Asset.Image == nil {
		err := fmt.Errorf("Failed to composite piece, no canvas on which to draw.")
		logErr.Println(err)
		return err
	}
	// Check for composition tree trunk.
	if len(p.Regions) == 0 {
		err := fmt.Errorf("Failed to composit piece, no tree to composite.")
	}
	// Find non-nil branches.
	var anyRegion bool
	for _, region := range p.Regions {
		// No region here, move on.
		if region == nil {
			continue
		}
		anyRegion = true
		// Climb the tree!
		comp, err := region.Composite()
		if err != nil {
			err := fmt.Errorf("Failed to composite piece: %s", err)
			logErr.Println(err)
			return err
		}
		// Expand the current canvas if necessary.
		if !comp.Bounds().In(p.Asset.Image.Bounds()) {
			// Replace the old canvas with the new one.
			p.Asset.Image = GrowImage(p.Asset.Image, comp.Bounds())
		}
		// Composite onto canvas. @TODO: Include size here
		offset := CenterOffset(region.Coordinates(), comp.Bounds())
		draw.Over.Draw(p.Asset.Image, p.Asset.Image.Bounds(), comp, offset)
	}
	// Check that we had at least one branch to climb.
	if !anyRegion {
		err := fmt.Errorf("Failed to composite piece, no regions were initialized.")
		logErr.Println(err)
		return err

	}
	// We successfully composited every branch.
	// Now we composite the result with the piece.
	// @TODO: final composite
	log.Println("Composited Piece, %v", p)
	return nil
}
