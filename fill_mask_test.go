package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo/v2"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalFillMaskRequest(t *testing.T) {
	// No options
	{
		fmExpected := hfapigo.SummarizationRequest{
			Inputs: []string{"Please to be fill in this [MASK]"},
		}

		jsonBuf, err := json.Marshal(fmExpected)
		if err != nil {
			t.Fatal(err)
		}

		fmActual := hfapigo.SummarizationRequest{}
		err = json.Unmarshal(jsonBuf, &fmActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(fmExpected, fmActual) {
			t.Fatalf("Expected %v, got %v", fmExpected, fmActual)
		}
	}

	// Options
	{
		fmExpected := hfapigo.SummarizationRequest{
			Inputs:  []string{"Please fill in this [MASK]"},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		}

		jsonBuf, err := json.Marshal(fmExpected)
		if err != nil {
			t.Fatal(err)
		}

		fmActual := hfapigo.SummarizationRequest{}
		err = json.Unmarshal(jsonBuf, &fmActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(fmExpected, fmActual) {
			t.Fatalf("Expected %v, got %v", fmExpected, fmActual)
		}
	}
}

func TestFillMaskRequest(t *testing.T) {
	// Basic request
	{
		inputs := []string{"Please fill in this [MASK]"}

		fmresps, err := hfapigo.SendFillMaskRequest(hfapigo.RecommendedFillMaskModel, &hfapigo.FillMaskRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(fmresps) != len(inputs) {
			t.Fatalf("Expected %d responses, got %d", len(inputs), len(fmresps))
		}
		for _, resp := range fmresps {
			if len(resp.Masks) == 0 {
				t.Fatalf("Expected nonzero masks")
			}
			for _, mask := range resp.Masks {
				if mask.Sequence == "" {
					t.Fatal("Expected non-empty mask sequence")
				}
				if mask.Score == 0.0 {
					t.Fatal("Expected non-zero score")
				}
				if mask.TokenID == 0 {
					t.Fatal("Expected non-zero token ID")
				}
				if mask.TokenStr == "" {
					t.Fatal("Expected non-empty token string")
				}
			}
		}
	}

	// Multiple inputs
	{
		inputs := []string{
			"Please fill in this [MASK]",
			"Please fill in this [MASK] too",
		}

		fmresps, err := hfapigo.SendFillMaskRequest(hfapigo.RecommendedFillMaskModel, &hfapigo.FillMaskRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(fmresps) != len(inputs) {
			t.Fatalf("Expected %d responses, got %d", len(inputs), len(fmresps))
		}
		for _, resp := range fmresps {
			if len(resp.Masks) == 0 {
				t.Fatalf("Expected nonzero masks")
			}
			for _, mask := range resp.Masks {
				if mask.Sequence == "" {
					t.Fatal("Expected non-empty mask sequence")
				}
				if mask.Score == 0.0 {
					t.Fatal("Expected non-zero score")
				}
				if mask.TokenID == 0 {
					t.Fatal("Expected non-zero token ID")
				}
				if mask.TokenStr == "" {
					t.Fatal("Expected non-empty token string")
				}
			}
		}
	}

	// Multiple Masks
	{
		inputs := []string{
			"Please fill in this [MASK] as well as this [MASK]",
		}

		fmresps, err := hfapigo.SendFillMaskRequest(hfapigo.RecommendedFillMaskModel, &hfapigo.FillMaskRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(fmresps) != 2 {
			t.Fatalf("Expected %d response, got %d", 2, len(fmresps))
		}
		for _, resp := range fmresps {
			if len(resp.Masks) == 0 {
				t.Fatalf("Expected nonzero masks")
			}
			for _, mask := range resp.Masks {
				if mask.Sequence == "" {
					t.Fatal("Expected non-empty mask sequence")
				}
				if mask.Score == 0.0 {
					t.Fatal("Expected non-zero score")
				}
				if mask.TokenID == 0 {
					t.Fatal("Expected non-zero token ID")
				}
				if mask.TokenStr == "" {
					t.Fatal("Expected non-empty token string")
				}
			}
		}
	}

	// Invalid request
	{
		fmresps, err := hfapigo.SendFillMaskRequest(hfapigo.RecommendedFillMaskModel, &hfapigo.FillMaskRequest{
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatalf("expected error - invalid request")
		}
		if fmresps != nil {
			t.Fatalf("expected nil response - invalid request")
		}
	}
}
