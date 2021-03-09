package wistia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultBaseURL   = "https://api.wistia.st/v1/"
	defaultUserAgent = "wistia-go-client/1.0"
)

type Client struct {
	accessToken string // TODO: use oauth2 package?
	httpClient  *http.Client
	baseURL     string

	Projects *ProjectsProvider
}

type provider struct {
	client *Client
}

func NewClient(httpClient *http.Client, accessToken string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		accessToken: accessToken,
		httpClient:  httpClient,
		baseURL:     defaultBaseURL,
	}
	client.Projects = &ProjectsProvider{client}
	return client
}

func (c *Client) request(ctx context.Context, method, url string, body interface{}, responseType interface{}) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal json for project: %s", err)
	}
	log.Printf("[TRACE] Request body: %s", string(payload))

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	if body != nil {
		req.Header.Set("Content-type", "application/json")
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	log.Printf("[TRACE] API response: %v", resp)

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("failed to read response body: %s", err)
		}
		return resp, fmt.Errorf("the Wistia API responded with status %d and body %s", resp.StatusCode, string(respBody))
	}

	err = json.NewDecoder(resp.Body).Decode(responseType)
	if err != nil {
		return resp, fmt.Errorf("failed to decode JSON from response body: %s", err)
	}

	return resp, nil
}
