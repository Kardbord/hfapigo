package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnmarshalSentenceSimilarityRequest(t *testing.T) {
	// No options
	{
		expected := hfapigo.SentenceSimilarityRequest{
			Inputs: hfapigo.SentenceSimilarityInputs{
				SourceSentence: "That is a happy person",
				Sentences: []string{
					"That is a happy dog",
					"That is a very happy person",
					"Today is a sunny day",
				},
			},
		}

		jsonBuf, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}

		actual := hfapigo.SentenceSimilarityRequest{}
		err = json.Unmarshal(jsonBuf, &actual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(expected, actual) {
			t.Fatalf("Expected %v, got %v", expected, actual)
		}
	}

	// Options
	{
		{
			expected := hfapigo.SentenceSimilarityRequest{
				Inputs: hfapigo.SentenceSimilarityInputs{
					SourceSentence: "That is a happy person",
					Sentences: []string{
						"That is a happy dog",
						"That is a very happy person",
						"Today is a sunny day",
					},
				},
				Options: *hfapigo.NewOptions().SetWaitForModel(true).SetUseGPU(false),
			}

			jsonBuf, err := json.Marshal(expected)
			if err != nil {
				t.Fatal(err)
			}

			actual := hfapigo.SentenceSimilarityRequest{}
			err = json.Unmarshal(jsonBuf, &actual)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(expected, actual) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		}
	}
}

func TestSentenceSimilarityRequest(t *testing.T) {
	// Minimal Request
	{
		inputSentences := []string{"That is a happy dog"}
		resp, err := hfapigo.SendSentenceSimilarityRequest(hfapigo.RecommendedSentenceSimilarityModel, &hfapigo.SentenceSimilarityRequest{
			Inputs: hfapigo.SentenceSimilarityInputs{
				SourceSentence: "That is a happy person",
				Sentences:      inputSentences,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("unexpected nil SentenceSimilarityResponse")
		}
		if len(*resp) != len(inputSentences) {
			t.Fatalf("Expected %d number of responses, got %d", len(inputSentences), len(*resp))
		}
	}

	// Multiple Sentences
	{
		inputSentences := []string{"That is a happy dog", "That is a very happy person", "Today is a sunny day"}
		resp, err := hfapigo.SendSentenceSimilarityRequest(hfapigo.RecommendedSentenceSimilarityModel, &hfapigo.SentenceSimilarityRequest{
			Inputs: hfapigo.SentenceSimilarityInputs{
				SourceSentence: "That is a happy person",
				Sentences:      inputSentences,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("unexpected nil SentenceSimilarityResponse")
		}
		if len(*resp) != len(inputSentences) {
			t.Fatalf("Expected %d number of responses, got %d", len(inputSentences), len(*resp))
		}
	}

	// Invalid request
	{
		resp, err := hfapigo.SendSentenceSimilarityRequest(hfapigo.RecommendedSentenceSimilarityModel, &hfapigo.SentenceSimilarityRequest{
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatal("Expected nil response")
		}
	}
}
