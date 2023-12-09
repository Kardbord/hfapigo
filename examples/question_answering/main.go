package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Kardbord/hfapigo"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

func main() {
	context := `The Amazon rainforest (Portuguese: Floresta Amazônica or Amazônia; Spanish: Selva Amazónica, Amazonía or usually Amazonia; French: Forêt amazonienne; Dutch: Amazoneregenwoud), also known in English as Amazonia or the Amazon Jungle, is a moist broadleaf forest that covers most of the Amazon basin of South America. This basin encompasses 7,000,000 square kilometres (2,700,000 sq mi), of which 5,500,000 square kilometres (2,100,000 sq mi) are covered by the rainforest. This region includes territory belonging to nine nations. The majority of the forest is contained within Brazil, with 60% of the rainforest, followed by Peru with 13%, Colombia with 10%, and with minor amounts in Venezuela, Ecuador, Bolivia, Guyana, Suriname and French Guiana. States or departments in four nations contain "Amazonas" in their names. The Amazon represents over half of the planet's remaining rainforests, and comprises the largest and most biodiverse tract of tropical rainforest in the world, with an estimated 390 billion individual trees divided into 16,000 species.`
	question := "How many square miles is the Amazon Rainforest?"

	fmt.Printf("Context:\n%s\n\n", context)
	fmt.Printf("Question: %s\n\n", question)

	fmt.Print("Sending request")

	type ChanRv struct {
		resp *hfapigo.QuestionAnsweringResponse
		err  error
	}
	ch := make(chan ChanRv)

	go func() {
		qaResp, err := hfapigo.SendQuestionAnsweringRequest(hfapigo.RecommendedQuestionAnsweringModel, &hfapigo.QuestionAnsweringRequest{
			Inputs: hfapigo.QuestionAnsweringInputs{
				Question: question,
				Context:  context,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{qaResp, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			fmt.Println("\nAnswer:", chrv.resp.Answer)
			fmt.Printf("Confidence: %.2f%%\n", chrv.resp.Score*100)
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
