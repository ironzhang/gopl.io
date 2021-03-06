package thumbnail

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Image returns a thumbnail-size version of src.
func Image(src image.Image) image.Image {
	xs := src.Bounds().Size().X
	ys := src.Bounds().Size().Y
	width, heigth := 128, 128
	if aspect := float64(xs) / float64(ys); aspect < 1.0 {
		width = int(128 * aspect)
	} else {
		heigth = int(128 / aspect)
	}
	xscale := float64(xs) / float64(width)
	yscale := float64(ys) / float64(heigth)

	dst := image.NewRGBA(image.Rect(0, 0, width, heigth))

	// a very crude scaling algorithm
	for x := 0; x < width; x++ {
		for y := 0; y < heigth; y++ {
			srcx := int(float64(x) * xscale)
			srcy := int(float64(y) * yscale)
			dst.Set(x, y, src.At(srcx, srcy))
		}
	}
	return dst
}

func ImageStream(w io.Writer, r io.Reader) error {
	src, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	dst := Image(src)
	return jpeg.Encode(w, dst, nil)
}

func ImageFile2(outflie, infile string) (err error) {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outflie)
	if err != nil {
		return err
	}

	if err := ImageStream(out, in); err != nil {
		out.Close()
		return fmt.Errorf("scaling %s to %s: %v", infile, outflie, err)
	}
	return out.Close()
}

func ImageFile(infile string) (string, error) {
	ext := filepath.Ext(infile)
	outflie := strings.TrimSuffix(infile, ext) + ".thumb" + ext
	return outflie, ImageFile2(outflie, infile)
}
