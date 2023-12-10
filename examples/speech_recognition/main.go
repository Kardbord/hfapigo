package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Kardbord/hfapigo/v2"
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

	fmt.Printf("Sending speech recongition request (%s)", audioFile)

	type ChanRv struct {
		resp *hfapigo.SpeechRecognitionResponse
		err  error
	}
	ch := make(chan ChanRv)

	go func() {
		arResp := &hfapigo.SpeechRecognitionResponse{}
		var err error
		for i := 0; i < maxRetries; i++ {
			arResp, err = hfapigo.SendSpeechRecognitionRequest(hfapigo.RecommendedSpeechRecongnitionModelEnglish, audioFile)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 5)
			}
		}
		ch <- ChanRv{arResp, err}
	}()

	for {
		select {
		case chrv := <-ch:
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			fmt.Printf("\nRecognized text: \"%s\"\n", chrv.resp.Text)
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 200)
		}
	}

}
