package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

func main() {
	const audioFile = "./sample.flac"
	const maxRetries = 10

	fmt.Printf("Sending Audio Classification request (%s)", audioFile)

	type ChanRv struct {
		resps []*hfapigo.AudioClassificationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		acResps := []*hfapigo.AudioClassificationResponse{}
		var err error
		for i := 0; i < maxRetries; i++ {
			acResps, err = hfapigo.SendAudioClassificationRequest(hfapigo.RecommendedAudioClassificationModel, audioFile)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 5)
			}
		}
		ch <- ChanRv{acResps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			for _, r := range chrv.resps {
				fmt.Printf("Label: \"%s\" Score: %f\n", r.Label, r.Score)
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 200)
		}
	}

}
