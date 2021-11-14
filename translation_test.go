package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalTranslationRequest(t *testing.T) {
	// No options
	{
		trExpected := hfapigo.TranslationRequest{
			Input: "Меня зовут Вольфганг и я живу в Берлине",
		}

		jsonBuf, err := json.Marshal(trExpected)
		if err != nil {
			t.Fatal(err)
		}

		trActual := hfapigo.TranslationRequest{}
		err = json.Unmarshal(jsonBuf, &trActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(trExpected, trActual) {
			t.Fatalf("Expected %v, got %v", trExpected, trActual)
		}
	}

	// Options
	{
		trExpected := hfapigo.TranslationRequest{
			Input:   "Меня зовут Вольфганг и я живу в Берлине",
			Options: *hfapigo.NewOptions().SetWaitForModel(true).SetUseGPU(false),
		}

		jsonBuf, err := json.Marshal(trExpected)
		if err != nil {
			t.Fatal(err)
		}

		trActual := hfapigo.TranslationRequest{}
		err = json.Unmarshal(jsonBuf, &trActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(trExpected, trActual) {
			t.Fatalf("Expected %v, got %v", trExpected, trActual)
		}
	}
}

func TestTranslationRequest(t *testing.T) {
	// Minimal request
	{
		tresp, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Input: "Меня зовут Вольфганг и я живу в Берлине",
		})
		if err != nil {
			t.Fatal(err)
		}

		const expectedTranslationText = "My name is Wolfgang and I live in Berlin."
		if tresp.TranslationText != expectedTranslationText {
			t.Fatalf("Expected translation text %s, but got %s", expectedTranslationText, tresp.TranslationText)
		}
	}

	// Request with optional parameters
	{
		tresp, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Input:   "Меня зовут Вольфганг и я живу в Берлине",
			Options: *hfapigo.NewOptions().SetWaitForModel(true).SetUseGPU(false),
		})
		if err != nil {
			t.Fatal(err)
		}

		const expectedTranslationText = "My name is Wolfgang and I live in Berlin."
		if tresp.TranslationText != expectedTranslationText {
			t.Fatalf("Expected translation text %s, but got %s", expectedTranslationText, tresp.TranslationText)
		}
	}

	// Invalid request
	{
		tresp, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{})
		if err == nil {
			t.Fatalf("Expected an error - invalid request, got tresp=%v", *tresp)
		}
		if tresp != nil {
			t.Fatal("Expected nil response")
		}
	}
}
