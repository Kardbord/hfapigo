package main

import (
	"fmt"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

func main() {
	inputs := []string{
		"The tower is 324 metres (1,063 ft) tall, about the same height as an 81-storey building, and the tallest structure in Paris. Its base is square, measuring 125 metres (410 ft) on each side. During its construction, the Eiffel Tower surpassed the Washington Monument to become the tallest man-made structure in the world, a title it held for 41 years until the Chrysler Building in New York City was finished in 1930. It was the first structure to reach a height of 300 metres. Due to the addition of a broadcasting aerial at the top of the tower in 1957, it is now taller than the Chrysler Building by 5.2 metres (17 ft). Excluding transmitters, the Eiffel Tower is the second tallest free-standing structure in France after the Millau Viaduct.",
		"Along with Ford Prefect, Arthur Dent barely escapes from Earth as it is demolished to make way for a hyperspace bypass. Arthur spends the next several years, still wearing his dressing gown, helplessly launched from crisis to crisis while trying to straighten out his lifestyle. He rather enjoys tea, but seems to have trouble obtaining it in the far reaches of the galaxy. In time, he learns how to fly and carves a niche for himself as a sandwich-maker.",
	}

	fmt.Println("Inputs:")
	for _, input := range inputs {
		fmt.Println(input)
		fmt.Println()
	}
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.SummarizationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		sresps, err := hfapigo.SendSummarizationRequest(hfapigo.RecommmendedSummarizationModel, &hfapigo.SummarizationRequest{
			Inputs: inputs,
			Parameters: hfapigo.SummarizationParameters{
				DoSample: false,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{sresps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}

			for _, resp := range chrv.resps {
				fmt.Println("Summary:", resp.SummaryText)
			}

			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
