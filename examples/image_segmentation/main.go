package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
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
	fmt.Printf("Opening image: %s\n", inputImage)
	srcImg, err := OpenImg(inputImage)
	if err != nil {
		fmt.Println("Problem opening image:", err)
		return
	}

	fmt.Printf("Requesting segmentation of image: %s\n", inputImage)
	resps, err := hfapigo.SendImageSegmentationRequest(hfapigo.RecommendedImageSegmentationModel, inputImage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, r := range resps {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.Mask))
		mask, _, err := image.Decode(reader)
		if err != nil {
			fmt.Println("Problem decoding mask:", err)
			return
		}
		draw.DrawMask(srcImg, srcImg.Bounds(), mask, image.Point{}, mask, image.Point{}, draw.Src)
	}

	outf, err := os.Create(outputImage)
	if err != nil {
		fmt.Printf("Error creating %s -> %s\n", outputImage, err)
	}
	defer outf.Close()

	err = png.Encode(outf, srcImg)
	if err != nil {
		fmt.Println("Problem encoding output file:", err)
		return
	}
	fmt.Println("Output image written to", outf.Name())
}

func OpenImg(imgFile string) (draw.Image, error) {
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
