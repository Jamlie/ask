package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Question(apiKey string) func(context.Context, string) (*Response, error) {
	return func(ctx context.Context, s string) (*Response, error) {
		reqBody := Request{
			Contents: []contentRequest{
				{
					Parts: []partRequest{
						{
							Text: s,
						},
					},
				},
			},
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf(
				"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s",
				apiKey,
			),
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(res.Body)
			return nil, fmt.Errorf(
				"request failed with status %d: %s",
				res.StatusCode,
				string(bodyBytes),
			)
		}

		var response Response
		if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return &response, nil
	}
}
