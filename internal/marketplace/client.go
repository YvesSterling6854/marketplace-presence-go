package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	BaseURL, APIKey string
	HTTP            *http.Client
}
type Envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *APIError       `json:"error"`
}
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *Client) call(method, path string, body, out any) error {
	p, e := json.Marshal(body)
	if e != nil {
		return e
	}
	for a := 0; a < 4; a++ {
		r, e := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(p))
		if e != nil {
			return e
		}
		r.Header.Set("Authorization", "Bearer "+c.APIKey)
		r.Header.Set("Content-Type", "application/json")
		x, e := c.HTTP.Do(r)
		if e != nil {
			return e
		}
		raw, _ := io.ReadAll(x.Body)
		x.Body.Close()
		var v Envelope
		if e = json.Unmarshal(raw, &v); e != nil {
			return e
		}
		if !v.OK {
			if x.StatusCode == 429 && a < 3 {
				d := time.Duration(1<<a) * 200 * time.Millisecond
				if n, z := strconv.Atoi(x.Header.Get("Retry-After")); z == nil {
					d = time.Duration(n) * time.Second
				}
				time.Sleep(d)
				continue
			}
			if v.Error != nil {
				return fmt.Errorf("%s: %s", v.Error.Code, v.Error.Message)
			}
			return fmt.Errorf("request rejected")
		}
		if out != nil {
			return json.Unmarshal(v.Data, out)
		}
		return nil
	}
	return fmt.Errorf("rate limit retries exhausted")
}
func (c *Client) CreateChannel(ch, typ, vendor string) error {
	return c.call(http.MethodPost, "/v1/realtime/channel/create", map[string]any{"channel": ch, "type": typ, "vendor": vendor}, nil)
}
func (c *Client) IssueToken(id string, ch, cap []string, ttl int, out any) error {
	return c.call(http.MethodPost, "/v1/realtime/token/issue", map[string]any{"client_id": id, "channels": ch, "capabilities": cap, "ttl_seconds": ttl}, out)
}
func (c *Client) Publish(ch, event string, data any, account string) error {
	return c.call(http.MethodPost, "/v1/realtime/publish", map[string]any{"channel": ch, "event": event, "data": data, "account_id": account}, nil)
}
func (c *Client) Presence(ch string, out any) error {
	return c.call(http.MethodGet, "/v1/realtime/presence/get/"+ch, nil, out)
}
