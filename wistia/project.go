package wistia

import (
	"context"
	"fmt"
	"net/http"
)

type ProjectsProvider provider

type Project struct {
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	MediaCount                   int    `json:"media_count""`
	Created                      string `json:"created"`
	Updated                      string `json:"updated"`
	HashedId                     string `json:"hashedId"`
	AnonymousCanUpload           bool   `json:"anonymous_can_upload"`
	AnonymousCanUploadOldStyle   bool   `json:"anonymousCanUpload"` // XXX: Workaround for API bug
	AnonymousCanDownload         bool   `json:"anonymous_can_download"`
	AnonymousCanDownloadOldStyle bool   `json:"anonymousCanDownload"` // XXX: Workaround for API bug
	Public                       bool   `json:"public"`
	PublicId                     string `json:"public_id,omitempty"`
}

func (pp *ProjectsProvider) Create(ctx context.Context, p *Project) (*Project, error) {
	createdProject := &Project{}
	url := pp.client.baseURL + "projects.json"
	_, err := pp.client.request(ctx, http.MethodPost, url, p, createdProject)
	if err != nil {
		return nil, err
	}
	return createdProject, nil
}

func (pp *ProjectsProvider) Get(ctx context.Context, id string) (*Project, error) {
	project := &Project{}
	url := pp.client.baseURL + fmt.Sprintf("projects/%s.json", id)
	_, err := pp.client.request(ctx, http.MethodGet, url, nil, project)
	if err != nil {
		return nil, err
	}
	// XXX: Workaround for API bug
	project.AnonymousCanUpload = project.AnonymousCanUploadOldStyle
	project.AnonymousCanDownload = project.AnonymousCanDownloadOldStyle
	return project, nil
}

func (pp *ProjectsProvider) Update(ctx context.Context, p *Project) (*Project, error) {
	updatedProject := &Project{}
	url := pp.client.baseURL + fmt.Sprintf("projects/%s.json", p.HashedId)
	_, err := pp.client.request(ctx, http.MethodPut, url, p, updatedProject)
	if err != nil {
		return nil, err
	}
	return updatedProject, nil
}

func (pp *ProjectsProvider) Delete(ctx context.Context, p *Project) error {
	url := pp.client.baseURL + fmt.Sprintf("projects/%s.json", p.HashedId)
	if _, err := pp.client.request(ctx, http.MethodDelete, url, nil, nil); err != nil {
		return err
	}

	return nil
}
