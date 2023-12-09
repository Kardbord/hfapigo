package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalQARequest(t *testing.T) {
	// No options
	{
		qaExpected := hfapigo.QuestionAnsweringRequest{
			Inputs: hfapigo.QuestionAnsweringInputs{
				Question: "What's my name?",
				Context:  "My name is Clara and I live in Berkeley.",
			},
		}

		jsonBuf, err := json.Marshal(qaExpected)
		if err != nil {
			t.Fatal(err)
		}

		qaActual := hfapigo.QuestionAnsweringRequest{}
		err = json.Unmarshal(jsonBuf, &qaActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(qaExpected, qaActual) {
			t.Fatalf("Expected %v, got %v", qaExpected, qaActual)
		}
	}

	// Options
	{
		qaExpected := hfapigo.QuestionAnsweringRequest{
			Inputs: hfapigo.QuestionAnsweringInputs{
				Question: "What's my name?",
				Context:  "My name is Clara and I live in Berkeley.",
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		}

		jsonBuf, err := json.Marshal(qaExpected)
		if err != nil {
			t.Fatal(err)
		}

		qaActual := hfapigo.QuestionAnsweringRequest{}
		err = json.Unmarshal(jsonBuf, &qaActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(qaExpected, qaActual) {
			t.Fatalf("Expected %v, got %v", qaExpected, qaActual)
		}
	}
}

func TestQARequest(t *testing.T) {
	// Basic request
	{
		qaResp, err := hfapigo.SendQuestionAnsweringRequest(hfapigo.RecommendedQuestionAnsweringModel, &hfapigo.QuestionAnsweringRequest{
			Inputs: hfapigo.QuestionAnsweringInputs{
				Question: "What's my name?",
				Context:  "My name is Clara and I live in Berkeley.",
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if qaResp == nil {
			t.Fatalf("Expected non-nil response")
		}

		const (
			expectedAnswer = "Clara"
			expectedStart  = 11
			expectedEnd    = 16
		)
		if qaResp.Answer != expectedAnswer {
			t.Fatalf("Expected %s, got %s", expectedAnswer, qaResp.Answer)
		}
		if qaResp.Score == 0.0 {
			t.Fatal("Expected non-zero confidence")
		}
		if qaResp.Start != expectedStart {
			t.Fatalf("Expected %d, got %d", expectedStart, qaResp.Start)
		}
		if qaResp.End != expectedEnd {
			t.Fatalf("Expected %d, got %d", expectedEnd, qaResp.End)
		}
	}

	// Invalid request
	{
		qaResp, err := hfapigo.SendQuestionAnsweringRequest(hfapigo.RecommendedQuestionAnsweringModel, &hfapigo.QuestionAnsweringRequest{
			Inputs: hfapigo.QuestionAnsweringInputs{
				Question: "What's my name?",
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("Expected error - invalid request")
		}
		if qaResp != nil {
			t.Fatalf("Expected nil response - invalid request")
		}
	}
}
