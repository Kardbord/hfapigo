package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo/v2"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnMarshalConversationalRequest(t *testing.T) {
	// No options
	{
		crExpected := hfapigo.ConversationalRequest{
			Inputs: hfapigo.ConverstationalInputs{
				Text: "Hey my name is Julien! How are you?",
			},
		}

		jsonBuf, err := json.Marshal(crExpected)
		if err != nil {
			t.Fatal(err)
		}

		crActual := hfapigo.ConversationalRequest{}
		err = json.Unmarshal(jsonBuf, &crActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(crExpected, crActual) {
			t.Fatalf("Expected %v, got %v", crExpected, crActual)
		}
	}

	// Options
	{
		crExpected := hfapigo.ConversationalRequest{
			Inputs: hfapigo.ConverstationalInputs{
				Text: "Hey my name is Julien! How are you?",
			},
			Parameters: *(&hfapigo.ConversationalParameters{}).
				SetTempurature(0.2345).
				SetMinLength(10).
				SetMaxLength(20).
				SetRepetitionPenalty(20),
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		}

		jsonBuf, err := json.Marshal(crExpected)
		if err != nil {
			t.Fatal(err)
		}

		crActual := hfapigo.ConversationalRequest{}
		err = json.Unmarshal(jsonBuf, &crActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(crExpected, crActual) {
			t.Fatalf("Expected %v, got %v", crExpected, crActual)
		}
	}
}

func TestConversationalRequest(t *testing.T) {
	// Basic request
	{
		cresp, err := hfapigo.SendConversationalRequest(hfapigo.RecommendedConversationalModel, &hfapigo.ConversationalRequest{
			Inputs: hfapigo.ConverstationalInputs{
				Text: "Hey my name is Julien! How are you?",
			},
			Parameters: *(&hfapigo.ConversationalParameters{}).
				SetTempurature(0.2345).
				SetMinLength(10).
				SetMaxLength(20).
				SetRepetitionPenalty(20),
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if cresp == nil {
			t.Fatal("Expected non-nil response")
		}
		if cresp.GeneratedText == "" {
			t.Fatal("Expected generated text")
		}
		if len(cresp.Conversation.GeneratedResponses) == 0 {
			t.Fatal("Expected non-empty previous responses")
		}
		if len(cresp.Conversation.PastUserInputs) == 0 {
			t.Fatal("Expected non-empty previous inputs")
		}
		if len(cresp.Conversation.GeneratedResponses) != len(cresp.Conversation.PastUserInputs) {
			t.Fatalf("Expected same number of past response and past user inputs, got %d and %d", len(cresp.Conversation.GeneratedResponses), len(cresp.Conversation.PastUserInputs))
		}
	}

	// Invalid request
	{
		cresp, err := hfapigo.SendConversationalRequest(hfapigo.RecommendedConversationalModel, &hfapigo.ConversationalRequest{
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("Expected error - invalid request")
		}
		if cresp != nil {
			t.Fatal("Expected nil response - invalid request")
		}
	}
}
