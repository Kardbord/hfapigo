package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalTextGenerationRequest(t *testing.T) {
	// No options
	{
		tgExpected := hfapigo.TextGenerationRequest{
			Inputs: []string{"The answer to the universe is"},
		}

		jsonBuf, err := json.Marshal(tgExpected)
		if err != nil {
			t.Fatal(err)
		}

		tgActual := hfapigo.TextGenerationRequest{}
		err = json.Unmarshal(jsonBuf, &tgActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tgExpected, tgActual) {
			t.Fatalf("Expected %v, got %v", tgExpected, tgActual)
		}
	}

	// Options
	{
		tgExpected := hfapigo.TextGenerationRequest{
			Inputs: []string{"The answer to the universe is"},
			Parameters: *hfapigo.NewTextGenerationParameters().
				SetMaxTime(12.2).
				SetMaxNewTokens(240).
				SetReturnFullText(false),
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		}

		jsonBuf, err := json.Marshal(tgExpected)
		if err != nil {
			t.Fatal(err)
		}

		tgActual := hfapigo.TextGenerationRequest{}
		err = json.Unmarshal(jsonBuf, &tgActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tgExpected, tgActual) {
			t.Fatalf("Expected %v, got %v", tgExpected, tgActual)
		}
	}
}

func TestTextGenerationRequest(t *testing.T) {
	// Basic request
	{
		inputs := []string{"The answer to the universe is"}
		const returnSeqs = 1
		tgresps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tgresps) != len(inputs) {
			t.Fatalf("expected %d response", len(inputs))
		}
		for i := range inputs {
			if len(tgresps[i].GeneratedTexts) != returnSeqs {
				t.Fatalf("expected non-empty list of generated texts")
			}
			for j := 0; j < returnSeqs; j++ {
				if tgresps[i].GeneratedTexts[j] == "" {
					t.Fatal("expected non-empty generated text")
				}
			}
		}
	}

	// More complicated request
	{
		inputs := []string{
			"The answer to the universe is",
			"There once was a ship that put to sea",
		}
		const returnSeqs = 3
		tgresps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Inputs:     inputs,
			Parameters: *hfapigo.NewTextGenerationParameters().SetRepetitionPenaly(50.235).SetReturnFullText(false).SetNumReturnSequences(returnSeqs),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tgresps) != len(inputs) {
			t.Fatalf("expected %d responses", len(inputs))
		}
		for i := range inputs {
			if len(tgresps[i].GeneratedTexts) != returnSeqs {
				t.Fatalf("expected non-empty list of generated texts")
			}
			for j := 0; j < returnSeqs; j++ {
				if tgresps[i].GeneratedTexts[j] == "" {
					t.Fatal("expected non-empty generated text")
				}
			}
		}
	}

	// Invalid request
	{
		tgresps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Parameters: *hfapigo.NewTextGenerationParameters().SetRepetitionPenaly(50.235).SetReturnFullText(false),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("expected error - invalid request")
		}
		if tgresps != nil {
			t.Fatal("expected nil response - invalid request")
		}
	}
}
