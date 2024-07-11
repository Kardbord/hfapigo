package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo/v3"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalTranslationRequest(t *testing.T) {
	// No options
	{
		trExpected := hfapigo.TranslationRequest{
			Inputs: []string{"Меня зовут Вольфганг и я живу в Берлине"},
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
			Inputs:  []string{"Меня зовут Вольфганг и я живу в Берлине"},
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
		tresps, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Inputs:  []string{"Меня зовут Вольфганг и я живу в Берлине"},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tresps) == 0 {
			t.Fatal("Expected nonzero number of response objects")
		}

		const expectedTranslationText = "My name is Wolfgang and I live in Berlin."
		if tresps[0].TranslationText != expectedTranslationText {
			t.Fatalf("Expected translation text %s, but got %s", expectedTranslationText, tresps[0].TranslationText)
		}
	}

	// Multiple inputs
	{
		inputs := []string{
			"Меня зовут Вольфганг и я живу в Берлине",
			"Здравствуйте, не могли бы вы направить меня к автобусной остановке?",
		}
		tresps, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tresps) != len(inputs) {
			t.Fatalf("Expected %d number of responses, got %d", len(inputs), len(tresps))
		}

		const expectedTranslationText1 = "My name is Wolfgang and I live in Berlin."
		if tresps[0].TranslationText != expectedTranslationText1 {
			t.Fatalf("Expected translation text %s, but got %s", expectedTranslationText1, tresps[0].TranslationText)
		}

		const expectedTranslationText2 = "Hello, could you direct me to the bus stop?"
		if tresps[1].TranslationText != expectedTranslationText2 {
			t.Fatalf("Expected translation text %s, but got %s", expectedTranslationText2, tresps[1].TranslationText)
		}
	}

	// Request with optional parameters
	{
		tresps, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{
			Inputs:  []string{"Меня зовут Вольфганг и я живу в Берлине"},
			Options: *hfapigo.NewOptions().SetWaitForModel(true).SetUseGPU(false),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tresps) == 0 {
			t.Fatal("Expected nonzero number of response objects")
		}

		const expectedTranslationText = "My name is Wolfgang and I live in Berlin."
		if tresps[0].TranslationText != expectedTranslationText {
			t.Fatalf("Expected translation text %s, but got %s", expectedTranslationText, tresps[0].TranslationText)
		}
	}

	// Invalid request
	{
		tresps, err := hfapigo.SendTranslationRequest(hfapigo.RecommendedRussianToEnglishModel, &hfapigo.TranslationRequest{})
		if err == nil {
			t.Fatalf("Expected an error - invalid request, got tresp=%v", tresps)
		}
		if tresps != nil {
			t.Fatal("Expected nil response")
		}
	}
}
