package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalTQARequest(t *testing.T) {
	// No options
	{
		tqaExpected := hfapigo.TableQuestionAnsweringRequest{
			Inputs: hfapigo.TableQuestionAnsweringInputs{
				Query: "What is the population of Anytown?",
				Table: map[string][]string{
					"City":       {"Anytown"},
					"Population": {"12345"},
				},
			},
		}

		jsonBuf, err := json.Marshal(tqaExpected)
		if err != nil {
			t.Fatal(err)
		}

		tqaActual := hfapigo.TableQuestionAnsweringRequest{}
		err = json.Unmarshal(jsonBuf, &tqaActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tqaExpected, tqaActual) {
			t.Fatalf("Expected %v, got %v", tqaExpected, tqaActual)
		}
	}

	// Options
	{
		tqaExpected := hfapigo.TableQuestionAnsweringRequest{
			Inputs: hfapigo.TableQuestionAnsweringInputs{
				Query: "What is the population of Anytown?",
				Table: map[string][]string{
					"City":       {"Anytown", "Someplace"},
					"Population": {"12345", "7890"},
				},
			},
		}

		jsonBuf, err := json.Marshal(tqaExpected)
		if err != nil {
			t.Fatal(err)
		}

		tqaActual := hfapigo.TableQuestionAnsweringRequest{}
		err = json.Unmarshal(jsonBuf, &tqaActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(tqaExpected, tqaActual) {
			t.Fatalf("Expected %v, got %v", tqaExpected, tqaActual)
		}
	}
}

func TestTQARequest(t *testing.T) {
	// Basic request
	{
		tqaResp, err := hfapigo.SendTableQuestionAnsweringRequest(hfapigo.RecommendedTableQuestionAnsweringModel, &hfapigo.TableQuestionAnsweringRequest{
			Inputs: hfapigo.TableQuestionAnsweringInputs{
				Query: "What is the population of Anytown?",
				Table: map[string][]string{
					"City":       {"Anytown", "Someplace"},
					"Population": {"12345", "7890"},
				},
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if tqaResp == nil {
			t.Fatal("Expected non-nil response")
		}
		if tqaResp.Answer == "" {
			t.Fatal("Expected non empty answer")
		}
		if len(tqaResp.Coordinates) == 0 {
			t.Fatal("Expected non empty coordinates")
		}
		if len(tqaResp.Coordinates[0]) == 0 {
			t.Fatal("Expected non empty coordinates[0]")
		}
		if len(tqaResp.Cells) == 0 {
			t.Fatal("Expected non empty cells")
		}
		if tqaResp.Aggregator == "" {
			t.Fatal("Expected non empty aggregator")
		}
	}

	// Invalid request
	{
		tqaResp, err := hfapigo.SendTableQuestionAnsweringRequest(hfapigo.RecommendedTableQuestionAnsweringModel, &hfapigo.TableQuestionAnsweringRequest{
			Inputs: hfapigo.TableQuestionAnsweringInputs{
				Table: map[string][]string{
					"City":       {"Anytown", "Someplace"},
					"Population": {"12345", "7890"},
				},
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("Expected error - invalid request")
		}
		if tqaResp != nil {
			t.Fatal("Expected nil response - invalid request")
		}
	}
}
