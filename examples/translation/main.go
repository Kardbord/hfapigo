package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

func main() {
	inputs := []string{
		"Меня зовут Вольфганг и я живу в Берлине",
		"Здравствуйте, не могли бы вы направить меня к автобусной остановке?",
	}

	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.TranslationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		tresps, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Input:   inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
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

			fmt.Println()
			for _, resp := range chrv.resps {
				fmt.Println("Translation:", resp.TranslationText)
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}

}
