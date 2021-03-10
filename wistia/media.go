package wistia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type MediaProvider provider

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Asset struct {
	URL         string `json:"url"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	FileSize    int    `json:"fileSize"`
	ContentType string `json:"contentType"`
	Type        string `json:"type"`
}

type Media struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Project Project `json:"project"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
	Section string  `json:"section"`
	//Thumbnail   Thumbnail `json:"thumbnail"`
	Duration float64 `json:"duration"`
	Created  string  `json:"created"`
	Updated  string  `json:"updated"`
	//Assets      []Asset   `json:"assets"`
	//EmbedCode   string    `json:"embedCode"`
	Description string `json:"description"`
	HashedId    string `json:"hashed_id"`
}

func (mp *MediaProvider) CreateFromReader(ctx context.Context, m *Media, r io.Reader, filename string) (*Media, error) {
	req, err := mp.client.newRequest(ctx, http.MethodPost, mp.client.UploadBaseEndpoint)
	if err != nil {
		return nil, err
	}

	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	req.Body = pipeReader

	go func() {
		defer pipeWriter.Close()
		defer multipartWriter.Close()

		if err := multipartWriter.WriteField("name", m.Name); err != nil {
			fmt.Fprintf(os.Stderr, "error creating name field: %s", err)
			return
		}
		if err := multipartWriter.WriteField("description", m.Description); err != nil {
			fmt.Fprintf(os.Stderr, "error creating description field: %s", err)
			return
		}
		if err := multipartWriter.WriteField("project_id", m.Project.HashedId); err != nil {
			fmt.Fprintf(os.Stderr, "error creating project_id field: %s", err)
			return
		}
		if err := multipartWriter.WriteField("access_token", mp.client.accessToken); err != nil {
			fmt.Fprintf(os.Stderr, "error creating access token field: %s", err)
			return
		}
		formFile, err := multipartWriter.CreateFormFile("file", filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating form file: %s", err)
			return
		}
		if _, err := io.Copy(formFile, r); err != nil {
			fmt.Fprintf(os.Stderr, "error during copy: %s", err)
			return
		}
	}()

	resp, err := mp.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("the Wistia API responded with an error while creating the media: %s", body)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	createdMedia := &Media{}
	err = json.NewDecoder(resp.Body).Decode(createdMedia)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON from response body: %s", err)
	}
	return createdMedia, nil
}

func (mp *MediaProvider) CreateFromURL(ctx context.Context, m *Media, sourceAssetUrl string) (*Media, error) {
	req, err := mp.client.newRequest(ctx, http.MethodPost, mp.client.UploadBaseEndpoint)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	values := url.Values{
		"access_token": {mp.client.accessToken},
		"name":         {m.Name},
		"project_id":   {m.Project.HashedId},
		"url":          {sourceAssetUrl},
	}
	req.Body = ioutil.NopCloser(strings.NewReader(values.Encode()))

	resp, err := mp.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("the Wistia API responded with error %d while creating the media: %s", resp.StatusCode, body)
	}

	createdMedia := &Media{}
	err = json.NewDecoder(resp.Body).Decode(createdMedia)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON from response body: %s", err)
	}
	return createdMedia, nil
}

func (mp *MediaProvider) Get(ctx context.Context, id string) (*Media, error) {
	media := &Media{}
	url := mp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s.json", id)
	_, err := mp.client.request(ctx, http.MethodGet, url, nil, media)
	if err != nil {
		return nil, err
	}
	return media, nil
}

func (mp *MediaProvider) Update(ctx context.Context, m *Media) (*Media, error) {
	url := mp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s.json", m.HashedId)
	updatedMedia := &Media{}
	_, err := mp.client.request(ctx, http.MethodPut, url, m, updatedMedia)
	if err != nil {
		return nil, err
	}
	return updatedMedia, nil
}

func (mp *MediaProvider) Delete(ctx context.Context, m *Media) error {
	url := mp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s.json", m.HashedId)
	_, err := mp.client.request(ctx, http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
