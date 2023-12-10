package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo/v2"
	"github.com/google/go-cmp/cmp"
)

func TestMarhshalUnMarshalTextClassificationRequest(t *testing.T) {
	// No options
	{
		tcExpected := hfapigo.TextClassificationRequest{
			Inputs: []string{"You know, I find you quite fascinating."},
		}

		jsonBuf, err := json.Marshal(tcExpected)
		if err != nil {
			t.Fatal(err)
		}

		trActual := hfapigo.TextClassificationRequest{}
		err = json.Unmarshal(jsonBuf, &trActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tcExpected, trActual) {
			t.Fatalf("Expected %v, got %v", tcExpected, trActual)
		}
	}

	// Options
	{
		tcExpected := hfapigo.TextClassificationRequest{
			Inputs: []string{
				"You know, I find you quite fascinating.",
				"I don't really care for your disposition.",
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		}

		jsonBuf, err := json.Marshal(tcExpected)
		if err != nil {
			t.Fatal(err)
		}

		trActual := hfapigo.TextClassificationRequest{}
		err = json.Unmarshal(jsonBuf, &trActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tcExpected, trActual) {
			t.Fatalf("Expected %v, got %v", tcExpected, trActual)
		}
	}
}

func TestTextClassificationRequest(t *testing.T) {
	// Minimal request
	{
		tcresps, err := hfapigo.SendTextClassificationRequest(hfapigo.RecommendedTextClassificationModel, &hfapigo.TextClassificationRequest{
			Inputs:  []string{"You know, I find you quite fascinating."},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tcresps) != 1 {
			t.Fatal("Expected 1 response object")
		}
		if len(tcresps[0].Labels) == 0 {
			t.Fatal("Expected at least one label in response")
		}
		if tcresps[0].Labels[0].Name == "" {
			t.Fatal("Expected non-empty label name")
		}
		if tcresps[0].Labels[0].Score == 0.0 {
			t.Fatal("Expected nonzero score")
		}
	}

	// Multiple inputs
	{
		tcresps, err := hfapigo.SendTextClassificationRequest(hfapigo.RecommendedTextClassificationModel, &hfapigo.TextClassificationRequest{
			Inputs: []string{
				"You know, I find you quite fascinating.",
				"I don't really care for your disposition.",
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tcresps) != 2 {
			t.Fatal("Expected 2 response objects")
		}
		if len(tcresps[0].Labels) == 0 {
			t.Fatal("Expected at least one label in response")
		}
		if tcresps[0].Labels[0].Name == "" {
			t.Fatal("Expected non-empty label name")
		}
		if tcresps[0].Labels[0].Score == 0.0 {
			t.Fatal("Expected nonzero score")
		}
		if len(tcresps[1].Labels) == 0 {
			t.Fatal("Expected at least one label in response")
		}
		if tcresps[1].Labels[0].Name == "" {
			t.Fatal("Expected non-empty label")
		}
		if tcresps[1].Labels[0].Score == 0.0 {
			t.Fatal("Expected nonzero score")
		}
	}

	// Invalid request
	{
		tcresps, err := hfapigo.SendTextClassificationRequest(hfapigo.RecommendedTextClassificationModel, &hfapigo.TextClassificationRequest{
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("Expected error - invalid request")
		}
		if tcresps != nil {
			t.Fatal("Expected nil response - invalid request")
		}
	}
}
