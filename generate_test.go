package artwork

import (
	"image"
	_ "image/jpeg"
	"testing"
)

func TestComposite(t *testing.T) {
	// t.Fatal("not implemented")
	p := NewPiece(1, image.NewNRGBA(image.Pt(500, 500)), &image.Rectangle{1000, 1000})

	r := NewRegion()
	a := NewAsset()
	//a.Path = ""
	a.Load()
	a.Parent = r
	r.Asset = a
	p.Asset.Regions = append(p.Asset.Regions, r)

	p.Composite()
}
