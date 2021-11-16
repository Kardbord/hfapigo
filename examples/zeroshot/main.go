package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

func main() {
	inputs := []string{
		"Hi, I recently bought a device from your company but it is not working as advertised and I would like to get reimbursed!",
		"Hi, I am having a difficult time using your product. Do you have a manual or perhaps a list of frequently asked questions to help get me going?",
	}
	candidateLabels := []string{"refund", "legal", "faq"}
	model := hfapigo.RecommendedZeroShotModel

	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("CandidateLabels: %v\n", candidateLabels)
	fmt.Printf("Model: %s\n", hfapigo.RecommendedZeroShotModel)
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.ZeroShotResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		zresps, err := hfapigo.SendZeroShotRequest(model, &hfapigo.ZeroShotRequest{
			Inputs: inputs,
			Parameters: hfapigo.ZeroShotParameters{
				CandidateLabels: candidateLabels,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{zresps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}

			fmt.Println()
			for _, resp := range chrv.resps {
				fmt.Println("\nSequence:", resp.Sequence)
				for i, label := range resp.Labels {
					fmt.Printf("%s: %f\n", label, resp.Scores[i])
				}
			}
			return

		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}