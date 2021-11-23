package main

import (
	"fmt"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
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
	inputImg  = "./test-image.png"
	outputImg = "./test-image-output.png"

	detectRetries = 10
	detectBackoff = time.Second * 2
)

func main() {
	fmt.Printf("Opening image: %s\n", inputImg)
	img, err := OpenImg(inputImg)
	if err != nil {
		fmt.Println("Problem opening image:", err)
		return
	}

	objects, err := SendRequest(inputImg)
	if err != nil {
		fmt.Println("Problem during object detection request:", err)
		return
	}

	// Draw the output file
	for _, obj := range objects {
		col := color.RGBA{uint8(rand.Intn(266)), uint8(rand.Intn(266)), uint8(rand.Intn(266)), 255}
		Rect(obj.Box.XMin, obj.Box.YMin, obj.Box.XMax, obj.Box.YMax, img, col)
	}

	outf, err := os.Create(outputImg)
	if err != nil {
		fmt.Println("Problem creating output file:", err)
		return
	}
	defer outf.Close()

	err = png.Encode(outf, img)
	if err != nil {
		fmt.Println("Problem encoding output file:", err)
		return
	}
	fmt.Println("Output image written to", outf.Name())
}

func OpenImg(imgFile string) (draw.Image, error) {
	f, err := os.Open(inputImg)
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

func SendRequest(imgFile string) ([]*hfapigo.ObjectDetectionResponse, error) {
	fmt.Printf("Sending object detection request for image (%s)", imgFile)

	type ChanRv struct {
		resps []*hfapigo.ObjectDetectionResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		objects, err := DetectObjects(inputImg, detectRetries, detectBackoff)
		ch <- ChanRv{objects, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			return chrv.resps, chrv.err
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 400)
		}
	}
}

func DetectObjects(imgFile string, retries int, backoff time.Duration) ([]*hfapigo.ObjectDetectionResponse, error) {
	objects := []*hfapigo.ObjectDetectionResponse{}
	var err error
	for i := 0; i < retries; i++ {
		objects, err = hfapigo.SendObjectDetectionRequest(hfapigo.RecommendedObjectDetectionModel, inputImg)
		if err == nil {
			break
		}
		time.Sleep(backoff)
	}
	return objects, err
}

// HLine draws a horizontal line
func HLine(x1, y, x2 int, img draw.Image, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func VLine(x, y1, y2 int, img draw.Image, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(x1, y1, x2, y2 int, img draw.Image, col color.Color) {
	HLine(x1, y1, x2, img, col)
	HLine(x1, y2, x2, img, col)
	VLine(x1, y1, y2, img, col)
	VLine(x2, y1, y2, img, col)
}
