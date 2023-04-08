package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	rand.Seed(time.Now().UnixNano())
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

const (
	inputImage  = "./test-image.png"
	outputImage = "./test-image-output.png"
)

func main() {
	resps, err := SendRequest()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	srcImg, err := OpenImg(inputImage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = DrawMasks(srcImg, resps)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = WriteOutputFile(srcImg, outputImage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func OpenImg(imgFile string) (draw.Image, error) {
	fmt.Printf("Opening image: %s\n", imgFile)
	f, err := os.Open(inputImage)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	dimg, ok := img.(draw.Image)
	if !ok {
		return nil, fmt.Errorf("%T is not a drawable image type", img)
	}
	return dimg, err
}

func SendRequest() ([]*hfapigo.ImageSegmentationResponse, error) {
	fmt.Printf("Requesting segmentation of image: %s\n", inputImage)
	return hfapigo.SendImageSegmentationRequest(hfapigo.RecommendedImageSegmentationModel, inputImage)
}

// See https://go.dev/blog/image-draw
func DrawMasks(srcImg draw.Image, resps []*hfapigo.ImageSegmentationResponse) error {
	type LabelCoords struct {
		X int
		Y int
	}

	labelCoords := []LabelCoords{}
	for segmentN, r := range resps {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.Mask))
		mask, _, err := image.Decode(reader)
		if err != nil {
			return err
		}
		WriteOutputFile(mask, fmt.Sprintf("mask%d.png", segmentN))
		fmt.Printf("Label: %s, Score: %f\n", r.Label, r.Score)

		segment := image.NewRGBA(mask.Bounds())
		segmentColor := color.RGBA{
			R: uint8(rand.Intn(266)),
			G: uint8(rand.Intn(266)),
			B: uint8(rand.Intn(266)),
			A: 210,
		}

		// There's probably a better way to do this, but image
		// processing is not my forte.
		for x := 0; x < segment.Bounds().Dx(); x++ {
			for y := 0; y < segment.Bounds().Dy(); y++ {
				mr, mg, mb, ma := mask.At(x, y).RGBA()
				if mr+mg+mb == ma*3 {
					segment.SetRGBA(x, y, segmentColor)
					if len(labelCoords) == segmentN {
						labelCoords = append(labelCoords, LabelCoords{
							X: x + 5,
							Y: y + 20,
						})
					}
				}
			}
		}

		WriteOutputFile(segment, fmt.Sprintf("segment%d.png", segmentN))
		draw.DrawMask(srcImg, srcImg.Bounds(), segment, image.Point{}, mask, image.Point{}, draw.Over)
	}

	// Add labels last so they show through.
	for i := range resps {
		AddLabel(srcImg, labelCoords[i].X, labelCoords[i].Y, fmt.Sprintf("%s (%.2f%%)", resps[i].Label, resps[i].Score*100))
	}
	return nil
}

func WriteOutputFile(img image.Image, filename string) error {
	outf, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outf.Close()

	err = png.Encode(outf, img)
	if err != nil {
		return err
	}
	fmt.Println("Output image written to", outf.Name())
	return nil
}

func AddLabel(img draw.Image, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
