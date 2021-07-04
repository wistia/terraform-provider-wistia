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
	defaultAPIEndpoint    = "https://api.wistia.com/v1/"
	defaultUploadEndpoint = "https://upload.wistia.com/"
	defaultUserAgent      = "wistia-go-client/1.0"
)

type Client struct {
	accessToken string // TODO: use oauth2 package?
	httpClient  *http.Client

	APIBaseEndpoint    string
	UploadBaseEndpoint string

	Media          *MediaProvider
	Projects       *ProjectsProvider
	Customizations *CustomizationsProvider
}

type provider struct {
	client *Client
}

func NewClient(httpClient *http.Client, accessToken string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		accessToken:        accessToken,
		httpClient:         httpClient,
		APIBaseEndpoint:    defaultAPIEndpoint,
		UploadBaseEndpoint: defaultUploadEndpoint,
	}
	client.Media = &MediaProvider{client}
	client.Projects = &ProjectsProvider{client}
	client.Customizations = &CustomizationsProvider{client}
	return client
}

func (c *Client) request(ctx context.Context, method, url string, body interface{}, responseType interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, url)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req.WithContext(ctx), body, responseType)
}

func (c *Client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("User-Agent", defaultUserAgent)
	return req, nil
}

func (c *Client) doRequest(req *http.Request, body interface{}, responseType interface{}) (*http.Response, error) {
	if body != nil {
		req.Header.Set("Content-type", "application/json")
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("couldn't marshal json for body: %s", err)
		}
		log.Printf("[TRACE] Request body: %s", string(payload))
		req.Body = ioutil.NopCloser(bytes.NewReader(payload))
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	log.Printf("[TRACE] API response: %v", resp)

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("failed to read response body: %s", err)
		}
		return resp, fmt.Errorf("the Wistia API responded with status %d and body %s", resp.StatusCode, string(respBody))
	}

	if responseType != nil {
		err = json.NewDecoder(resp.Body).Decode(responseType)
		if err != nil {
			return resp, fmt.Errorf("failed to decode JSON from response body: %s", err)
		}
	}

	return resp, nil
}
