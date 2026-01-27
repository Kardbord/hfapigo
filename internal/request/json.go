package request

import (
	"encoding/json"
	"fmt"
	"io"
)

func DoJSON[TReq any, TResp any](
	opts RequestOptions,
	method string,
	path string,
	reqBody TReq,
) (TResp, error) {

	var zero TResp

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return zero, err
	}

	resp, err := DoBytes(
		opts,
		method,
		path,
		buf,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return zero, fmt.Errorf("hf api error (%d): %s", resp.StatusCode, string(b))
	}

	var out TResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return zero, err
	}

	return out, nil
}
