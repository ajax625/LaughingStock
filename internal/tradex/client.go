package tradex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type CreateSubscriptionRequest struct {
	Name          string  `json:"name"`
	WebhookURL    string  `json:"webhook_url"`
	Symbol        string  `json:"symbol"`
	TraderType    string  `json:"trader_type"`
	MinROI        float64 `json:"min_roi"`
	MinRR         float64 `json:"min_rr"`
	MinConfidence float64 `json:"min_confidence"`
}

type CreateSubscriptionResponse struct {
	ID string `json:"id"`
}

func (c *Client) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (string, error) {
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/subscriptions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("tradex: create subscription failed with status %d", resp.StatusCode)
	}

	var out CreateSubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.ID, nil
}

type UpdateSubscriptionRequest struct {
	WebhookURL    string  `json:"webhook_url,omitempty"`
	MinROI        float64 `json:"min_roi"`
	MinRR         float64 `json:"min_rr"`
	MinConfidence float64 `json:"min_confidence"`
}

func (c *Client) UpdateSubscription(ctx context.Context, id string, req UpdateSubscriptionRequest) error {
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+"/subscriptions/"+id, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tradex: update subscription failed with status %d", resp.StatusCode)
	}
	return nil
}

type Signal struct {
	ID          string   `json:"id"`
	Symbol      string   `json:"symbol"`
	Direction   string   `json:"direction"`
	Confidence  float64  `json:"confidence"`
	Entry       float64  `json:"entry"`
	Target      float64  `json:"target"`
	Stop        float64  `json:"stop"`
	Timeframe   string   `json:"timeframe"`
	SignalsFired []string `json:"signals_fired"`
	Mode        string   `json:"mode"`
	GeneratedAt string   `json:"generated_at"`
}

func (c *Client) GetLatestSignal(ctx context.Context, symbol string) (*Signal, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/signals?symbol=%s&limit=1", c.baseURL, symbol), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Signals []Signal `json:"signals"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Signals) == 0 {
		return nil, nil
	}
	return &out.Signals[0], nil
}

// ResearchJob is the response shape from TradeX's research job endpoints.
type ResearchJob struct {
	ID          string          `json:"id"`
	Symbol      string          `json:"symbol"`
	Status      string          `json:"status"`
	Result      json.RawMessage `json:"result"`
	ErrorMsg    string          `json:"error_msg,omitempty"`
	CreatedAt   string          `json:"created_at"`
	CompletedAt *string         `json:"completed_at,omitempty"`
}

// TriggerDeepResearch fires a fundamentals-only deep research job for the symbol.
// Returns the job_id. Idempotent — TradeX returns the existing job if one already
// completed today.
func (c *Client) TriggerDeepResearch(ctx context.Context, symbol string) (string, error) {
	body, _ := json.Marshal(map[string]string{"symbol": symbol})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/research/deep", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tradex: trigger research failed with status %d", resp.StatusCode)
	}

	var out struct {
		JobID string `json:"job_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.JobID, nil
}

// GetLatestResearch returns the most recent completed research job for the symbol.
func (c *Client) GetLatestResearch(ctx context.Context, symbol string) (*ResearchJob, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/research/jobs?symbol=%s&limit=1", c.baseURL, symbol), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Jobs []struct {
			ID          string  `json:"id"`
			Symbol      string  `json:"symbol"`
			Status      string  `json:"status"`
			HasResult   bool    `json:"has_result"`
			CreatedAt   string  `json:"created_at"`
			CompletedAt *string `json:"completed_at"`
		} `json:"jobs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Jobs) == 0 {
		return nil, nil
	}

	j := out.Jobs[0]
	if !j.HasResult {
		return &ResearchJob{ID: j.ID, Symbol: j.Symbol, Status: j.Status}, nil
	}

	// Fetch the full job with result
	jobReq, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/research/jobs/%s", c.baseURL, j.ID), nil)
	if err != nil {
		return nil, err
	}
	jobResp, err := c.httpClient.Do(jobReq)
	if err != nil {
		return nil, err
	}
	defer jobResp.Body.Close()

	var job ResearchJob
	if err := json.NewDecoder(jobResp.Body).Decode(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (c *Client) DeleteSubscription(ctx context.Context, id string) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/subscriptions/"+id, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tradex: delete subscription failed with status %d", resp.StatusCode)
	}
	return nil
}
