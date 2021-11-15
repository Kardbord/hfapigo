package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnmarshalSummarizationRequest(t *testing.T) {
	// No options
	{
		srExpected := hfapigo.SummarizationRequest{
			Inputs: []string{"Foobarbaz"},
		}

		jsonBuf, err := json.Marshal(srExpected)
		if err != nil {
			t.Fatal(err)
		}

		srActual := hfapigo.SummarizationRequest{}
		err = json.Unmarshal(jsonBuf, &srActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(srExpected, srActual) {
			t.Fatalf("Expected %v, got %v", srExpected, srActual)
		}
	}

	// Options
	{
		srExpected := hfapigo.SummarizationRequest{
			Inputs: []string{"Foobar", "baz"},
			Parameters: *(&hfapigo.SummarizationParameters{
				MaxLength:         5,
				TopK:              20,
				TopP:              1.25,
				RepetitionPenalty: 0.215,
				DoSample:          false,
			}).SetTempurature(92.123456789),
			Options: *hfapigo.NewOptions().SetUseCache(false),
		}

		jsonBuf, err := json.Marshal(srExpected)
		if err != nil {
			t.Fatal(err)
		}

		srActual := hfapigo.SummarizationRequest{}
		err = json.Unmarshal(jsonBuf, &srActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(srExpected, srActual) {
			t.Fatalf("Expected %v, got %v", srExpected, srActual)
		}
	}
}

func TestSummarizationRequest(t *testing.T) {

	// Basic request
	{
		inputs := []string{
			"The tower is 324 metres (1,063 ft) tall, about the same height as an 81-storey building, and the tallest structure in Paris. Its base is square, measuring 125 metres (410 ft) on each side. During its construction, the Eiffel Tower surpassed the Washington Monument to become the tallest man-made structure in the world, a title it held for 41 years until the Chrysler Building in New York City was finished in 1930. It was the first structure to reach a height of 300 metres. Due to the addition of a broadcasting aerial at the top of the tower in 1957, it is now taller than the Chrysler Building by 5.2 metres (17 ft). Excluding transmitters, the Eiffel Tower is the second tallest free-standing structure in France after the Millau Viaduct.",
			"Along with Ford Prefect, Arthur Dent barely escapes from Earth as it is demolished to make way for a hyperspace bypass. Arthur spends the next several years, still wearing his dressing gown, helplessly launched from crisis to crisis while trying to straighten out his lifestyle. He rather enjoys tea, but seems to have trouble obtaining it in the far reaches of the galaxy. In time, he learns how to fly and carves a niche for himself as a sandwich-maker.",
		}

		sresps, err := hfapigo.SendSummarizationRequest(hfapigo.RecommmendedSummarizationModel, &hfapigo.SummarizationRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(sresps) != len(inputs) {
			t.Fatalf("Expected %d number of responses, got %d", len(inputs), len(sresps))
		}
	}

	// Invalid request
	{
		sresps, err := hfapigo.SendSummarizationRequest(hfapigo.RecommmendedSummarizationModel, &hfapigo.SummarizationRequest{
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatalf("expected error - invalid request")
		}
		if sresps != nil {
			t.Fatalf("expected nil response - invalid request")
		}
	}
}
