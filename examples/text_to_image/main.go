package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime"
	"os"
	"time"

	"github.com/Kardbord/hfapigo/v3"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

func main() {
	fmt.Print("Enter an image prompt: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print(err)
		return
	}

	type ChanRv struct {
		resp   image.Image
		format string
		err    error
	}
	ch := make(chan ChanRv)

	fmt.Print("Sending request")
	go func() {
		img, fmt, err := hfapigo.SendTextToImageRequest(hfapigo.RecommendedTextToImageModel, &hfapigo.TextToImageRequest{
			Inputs:  input,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		ch <- ChanRv{img, fmt, err}
	}()

	for {
		select {
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 300)
		case chrv := <-ch:
			if chrv.err != nil {
				fmt.Printf("\nError from Hugging Face: %s\n", chrv.err)
				return
			}

			filename := fmt.Sprintf("output.%s", chrv.format)
			fout, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer fout.Close()

			mimetype := mime.TypeByExtension(fmt.Sprintf(".%s", chrv.format))

			switch mimetype {
			case "image/jpeg":
				err = jpeg.Encode(fout, chrv.resp, nil)
			case "image/png":
				err = png.Encode(fout, chrv.resp)
			default:
				err = fmt.Errorf("unknown image format: %s", chrv.format)
			}

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\nWrote image to %s\n", filename)
			}

			return
		}
	}
}
