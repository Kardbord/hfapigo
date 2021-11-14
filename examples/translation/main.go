package main

import (
	"fmt"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

func main() {
	input := "Меня зовут Вольфганг и я живу в Берлине"

	fmt.Printf("Inputs: %s\n", input)
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resp *hfapigo.TranslationResponse
		err  error
	}
	ch := make(chan ChanRv)

	go func() {
		tresps, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Input: input,
		})

		ch <- ChanRv{tresps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			fmt.Println("\nTranslation:", chrv.resp.TranslationText)
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 50)
		}
	}

}
