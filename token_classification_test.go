package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalTokenClassificationRequest(t *testing.T) {
	// No options
	{
		tcExpected := hfapigo.TokenClassificationRequest{
			Inputs: []string{"My name is Sarah Jessica Parker but you can call me Jessica"},
		}

		jsonBuf, err := json.Marshal(tcExpected)
		if err != nil {
			t.Fatal(err)
		}

		trActual := hfapigo.TokenClassificationRequest{}
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
		tcExpected := hfapigo.TokenClassificationRequest{
			Inputs:     []string{"My name is Sarah Jessica Parker but you can call me Jessica"},
			Parameters: *hfapigo.NewTokenClassificationParameters().SetAggregationStrategy(hfapigo.AggregationStrategyFirst),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
		}

		jsonBuf, err := json.Marshal(tcExpected)
		if err != nil {
			t.Fatal(err)
		}

		trActual := hfapigo.TokenClassificationRequest{}
		err = json.Unmarshal(jsonBuf, &trActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tcExpected, trActual) {
			t.Fatalf("Expected %v, got %v", tcExpected, trActual)
		}
	}
}

func TestTokenClassificationRequest(t *testing.T) {
	// Minimal request
	{
		tcresps, err := hfapigo.SendTokenClassificationRequest(hfapigo.RecommendedTokenClassificationModel, &hfapigo.TokenClassificationRequest{
			Inputs:  []string{"My name is Sarah Jessica Parker but you can call me Jessica"},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tcresps) != 1 {
			t.Fatal("expected 1 response")
		}
		for _, r := range tcresps {
			for _, g := range r.Entities {
				if g.Label == "" {
					t.Fatal("expected non-empty entity list")
				}
				if g.Score == 0.0 {
					t.Fatal("Expected non-zero score")
				}
				if g.Entity == "" {
					t.Fatal("Expected non-empty entity")
				}
				if g.Start == 0 {
					t.Fatal("Expected non-zero start")
				}
				if g.End == 0 {
					t.Fatal("Expected non-zero end")
				}
			}
		}
	}

	// Multiple inputs and parameters
	{
		tcresps, err := hfapigo.SendTokenClassificationRequest(hfapigo.RecommendedTokenClassificationModel, &hfapigo.TokenClassificationRequest{
			Inputs: []string{
				"My name is Sarah Jessica Parker but you can call me Jessica",
				"My name is Clara and I live in Berkeley, California.",
			},
			Parameters: *hfapigo.NewTokenClassificationParameters().SetAggregationStrategy(hfapigo.AggregationStrategyNone),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(tcresps) != 2 {
			t.Fatal("expected 2 responses")
		}
		for _, r := range tcresps {
			for _, g := range r.Entities {
				if g.Label == "" {
					t.Fatal("expected non-empty entity list")
				}
				if g.Score == 0.0 {
					t.Fatal("Expected non-zero score")
				}
				if g.Entity == "" {
					t.Fatal("Expected non-empty entity")
				}
				if g.Start == 0 {
					t.Fatal("Expected non-zero start")
				}
				if g.End == 0 {
					t.Fatal("Expected non-zero end")
				}
			}
		}
	}

	// Invalid request
	{
		tcresps, err := hfapigo.SendTokenClassificationRequest(hfapigo.RecommendedTokenClassificationModel, &hfapigo.TokenClassificationRequest{
			Parameters: *hfapigo.NewTokenClassificationParameters().SetAggregationStrategy(hfapigo.AggregationStrategyNone),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("Expected error - invalid request")
		}
		if tcresps != nil {
			t.Fatal("Expected nil response - invalid request")
		}
	}
}
