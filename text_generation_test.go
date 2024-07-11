package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo/v2"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalTextGenerationRequest(t *testing.T) {
	// No options
	{
		tgExpected := hfapigo.TextGenerationRequest{
			Input: "The answer to the universe is",
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
			Input: "The answer to the universe is",
			Parameters: *hfapigo.NewTextGenerationParameters().
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
		input := "The answer to the universe is"
		tgresps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Input:   input,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tgresps) != 1 {
			t.Fatalf("expected 1 response, got %d", len(tgresps))
		}
		if tgresps[0].GeneratedText == "" {
			t.Fatal("expected non-empty generated text")
		}
	}

	// More complicated request
	{
		input := "There once was a ship that put to sea"
		tgresps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Input:      input,
			Parameters: *hfapigo.NewTextGenerationParameters().SetRepetitionPenaly(50.235).SetReturnFullText(false).SetDetails(true),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tgresps) != 1 {
			t.Fatalf("expected 1 response, got %d", len(tgresps))
		}
		if tgresps[0].GeneratedText == "" {
			t.Fatal("expected non-empty generated text")
		}
		if tgresps[0].Details.FinishReason == "" {
			t.Fatal("expected non-empty finish reason")
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
