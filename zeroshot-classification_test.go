package hfapigo_test

import (
	"encoding/json"
	"testing"

	"github.com/TannerKvarfordt/hfapigo"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalUnmarshalZeroShotRequest(t *testing.T) {
	// No options
	{
		zsrExpected := hfapigo.ZeroShotRequest{
			Inputs: []string{"Hi, I recently bought a device from your company but it is not working as advertised and I would like to get reimbursed!"},
			Parameters: hfapigo.ZeroShotParameters{
				CandidateLabels: []string{"refund", "legal", "faq"},
			},
		}

		jsonBuf, err := json.Marshal(zsrExpected)
		if err != nil {
			t.Fatal(err)
		}

		zsrActual := hfapigo.ZeroShotRequest{}
		err = json.Unmarshal(jsonBuf, &zsrActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(zsrExpected, zsrActual) {
			t.Fatalf("Expected %v, got %v", zsrExpected, zsrActual)
		}
	}

	// Options
	{
		zsrExpected := hfapigo.ZeroShotRequest{
			Inputs: []string{"Hi, I recently bought a device from your company but it is not working as advertised and I would like to get reimbursed!"},
			Parameters: *(&hfapigo.ZeroShotParameters{
				CandidateLabels: []string{"refund", "legal", "faq"},
			}).SetMultiLabel(true),
			Options: *hfapigo.NewOptions().SetWaitForModel(true).SetUseGPU(false),
		}

		jsonBuf, err := json.Marshal(zsrExpected)
		if err != nil {
			t.Fatal(err)
		}

		zsrActual := hfapigo.ZeroShotRequest{}
		err = json.Unmarshal(jsonBuf, &zsrActual)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(zsrExpected, zsrActual) {
			t.Fatalf("Expected %v, got %v", zsrExpected, zsrActual)
		}
	}
}

func TestZeroShotRequest(t *testing.T) {
	// Minimal request
	{
		zreq := hfapigo.ZeroShotRequest{
			Inputs: []string{"Hi, I recently bought a device from your company but it is not working as advertised and I would like to get reimbursed!"},
			Parameters: hfapigo.ZeroShotParameters{
				CandidateLabels: []string{"refund", "legal", "faq"},
			},
		}

		zresps, err := hfapigo.SendZeroShotRequest(&zreq, hfapigo.APIBaseURL+hfapigo.RecommendedZeroShotModel)
		if err != nil {
			t.Fatal(err)
		}
		if len(zresps) == 0 {
			t.Fatal("ZeroShotResponses should not be empty")
		}
		if len(zresps) != len(zreq.Inputs) {
			t.Fatalf("Expected len(zreq.Inputs) == len(zresps), but %d != %d", len(zreq.Inputs), len(zresps))
		}

		zresp := zresps[0]
		if !cmp.Equal(zresp.Sequence, zreq.Inputs[0]) {
			t.Fatalf("Expected zresp.Sequence to equal inputs, but %v != %v", zresp.Sequence, zreq.Inputs)
		}
		if len(zresp.Labels) != len(zreq.Parameters.CandidateLabels) {
			t.Fatalf("Expected len(zresp.Labels) to equal len(candidateLabels), but %d != %d", len(zresp.Labels), len(zreq.Parameters.CandidateLabels))
		}
		if len(zresp.Scores) != len(zresp.Labels) {
			t.Fatalf("expected len(zresp.Scores) == len(zresp.Labels), but %d != %d", len(zresp.Scores), len(zresp.Labels))
		}
	}

	// Request with optional parameters
	{
		zreq := hfapigo.ZeroShotRequest{
			Inputs: []string{
				"Hi, I recently bought a device from your company but it is not working as advertised and I would like to get reimbursed!",
				"To whom it may concern, I purchased an item from your storefront, but its behavior is not as anticipated and I will be expecting a refund!",
			},
			Parameters: *(&hfapigo.ZeroShotParameters{
				CandidateLabels: []string{"refund", "legal", "faq"},
			}).SetMultiLabel(true),
			Options: *hfapigo.NewOptions().SetWaitForModel(true).SetUseGPU(false),
		}

		zresps, err := hfapigo.SendZeroShotRequest(&zreq, hfapigo.APIBaseURL+hfapigo.RecommendedZeroShotModel)
		if err != nil {
			t.Fatal(err)
		}
		if len(zresps) == 0 {
			t.Fatal("ZeroShotResponses should not be empty")
		}
		if len(zresps) != len(zreq.Inputs) {
			t.Fatalf("Expected len(zreq.Inputs) == len(zresps), but %d != %d", len(zreq.Inputs), len(zresps))
		}

		zresp := zresps[0]
		if !cmp.Equal(zresp.Sequence, zreq.Inputs[0]) {
			t.Fatalf("Expected zresp.Sequence to equal inputs, but %v != %v", zresp.Sequence, zreq.Inputs[0])
		}
		if len(zresp.Labels) != len(zreq.Parameters.CandidateLabels) {
			t.Fatalf("Expected len(zresp.Labels) to equal len(candidateLabels), but %d != %d", len(zresp.Labels), len(zreq.Parameters.CandidateLabels))
		}
		if len(zresp.Scores) != len(zresp.Labels) {
			t.Fatalf("expected len(zresp.Scores) == len(zresp.Labels), but %d != %d", len(zresp.Scores), len(zresp.Labels))
		}

		zresp = zresps[1]
		if !cmp.Equal(zresp.Sequence, zreq.Inputs[1]) {
			t.Fatalf("Expected zresp.Sequence to equal inputs, but %v != %v", zresp.Sequence, zreq.Inputs[1])
		}
		if len(zresp.Labels) != len(zreq.Parameters.CandidateLabels) {
			t.Fatalf("Expected len(zresp.Labels) to equal len(candidateLabels), but %d != %d", len(zresp.Labels), len(zreq.Parameters.CandidateLabels))
		}
		if len(zresp.Scores) != len(zresp.Labels) {
			t.Fatalf("expected len(zresp.Scores) == len(zresp.Labels), but %d != %d", len(zresp.Scores), len(zresp.Labels))
		}
	}
}
