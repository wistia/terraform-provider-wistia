package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/wistia/terraform-provider-wistia/wistia"
	"log"
)

func projectResource() *schema.Resource {
	return &schema.Resource{
		Create: createProject,
		Read:   readProject,
		Update: updateProject,
		Delete: deleteProject,
		// TODO: Do we need this?
		//Exists: isProject,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Exposed by the API but not supported by projects#create
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"media_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hashed_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"anonymous_can_upload": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"anonymous_can_download": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"public_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createProject(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	p := projectFromResource(d)
	p, err := wc.Projects.Create(context.Background(), p)
	if err != nil {
		return fmt.Errorf("couldn't create Wistia project: %s", err)
	}

	log.Printf("[TRACE] Newly created project: %v", p)

	applyProjectFieldsToResource(p, d)

	return nil
}

func readProject(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	p, err := wc.Projects.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("couldn't get Wistia project: %s", err)
	}

	applyProjectFieldsToResource(p, d)

	return nil
}

func updateProject(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	p := projectFromResource(d)
	log.Printf("[TRACE] Project before update: %v", p)
	p, err := wc.Projects.Update(context.Background(), p)
	if err != nil {
		return fmt.Errorf("couldn't update Wistia project: %s", err)
	}

	log.Printf("[TRACE] Project after update: %v", p)

	applyProjectFieldsToResource(p, d)

	return nil
}

func deleteProject(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	p := projectFromResource(d)
	if err := wc.Projects.Delete(context.Background(), p); err != nil {
		return fmt.Errorf("couldn't delete Wistia project: %s", err)
	}

	return nil
}

// Private helpers

func applyProjectFieldsToResource(p *wistia.Project, d *schema.ResourceData) {
	d.SetId(p.HashedId)
	d.Set("name", p.Name)
	d.Set("description", p.Description)
	d.Set("media_count", p.MediaCount)
	d.Set("created", p.Created)
	d.Set("updated", p.Updated)
	d.Set("hashed_id", p.HashedId)
	d.Set("anonymous_can_upload", p.AnonymousCanUpload)
	d.Set("anonymous_can_download", p.AnonymousCanDownload)
	d.Set("public", p.Public)
	d.Set("public_id", p.PublicId)
}

func projectFromResource(d *schema.ResourceData) *wistia.Project {
	return &wistia.Project{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		MediaCount:                   d.Get("media_count").(int),
		Created:                      d.Get("created").(string),
		Updated:                      d.Get("updated").(string),
		HashedId:                     d.Get("hashed_id").(string),
		AnonymousCanDownload:         d.Get("anonymous_can_download").(bool),
		AnonymousCanDownloadOldStyle: d.Get("anonymous_can_download").(bool),
		AnonymousCanUpload:           d.Get("anonymous_can_upload").(bool),
		AnonymousCanUploadOldStyle:   d.Get("anonymous_can_upload").(bool),
		Public:                       d.Get("public").(bool),
		PublicId:                     d.Get("public_id").(string),
	}
}
