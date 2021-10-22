package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/wistia/terraform-provider-wistia/internal/wistia"
	"log"
	"os"
	"path"
)

func mediaResource() *schema.Resource {
	return &schema.Resource{
		Create: createMedia,
		Read:   readMedia,
		Update: updateMedia,
		Delete: deleteMedia,
		// TODO: Do we need this?
		//Exists: isMedia,

		Schema: map[string]*schema.Schema{
			"file": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file", "url"},
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file", "url"},
			},
			"media_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"section": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"thumbnail": {
			//	Type:     schema.TypeMap,
			//	Computed: true,
			//},
			"duration": {
				Type:     schema.TypeFloat,
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
			// TODO
			//"assets": {
			//	Type:     schema.TypeList,
			//	Elem:     wistia.Asset{},
			//	Computed: true,
			//},
			//"embed_code": {
			//	Type:       schema.TypeString,
			//	Computed:   true,
			//	Deprecated: "If you want to programmatically embed videos, follow the \"construct an embed code\" guide",
			//},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hashed_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createMedia(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	media := mediaFromResource(d)
	var err error
	if filePath, ok := d.GetOk("file"); ok {
		filePath := filePath.(string)
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("couldn't open file '%s': %s", filePath, err)
		}
		defer f.Close()
		media, err = wc.Media.CreateFromReader(context.Background(), media, f, path.Base(filePath))
		if err != nil {
			return fmt.Errorf("couldn't create media: %s", err)
		}
	} else {
		url := d.Get("url").(string)
		media, err = wc.Media.CreateFromURL(context.Background(), media, url)
		if err != nil {
			return fmt.Errorf("couldn't create media: %s", err)
		}
	}

	log.Printf("[TRACE] Newly created media: %v", media)

	applyMediaFieldsToResource(media, d)

	return nil
}

func readMedia(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	media, err := wc.Media.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("couldn't get media: %s", err)
	}

	log.Printf("[TRACE] Read media: %v", media)

	applyMediaFieldsToResource(media, d)

	return nil
}

func updateMedia(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	media := mediaFromResource(d)
	media, err := wc.Media.Update(context.Background(), media)
	if err != nil {
		return fmt.Errorf("couldn't update media: %s", err)
	}

	log.Printf("[TRACE] Read media: %v", media)

	applyMediaFieldsToResource(media, d)

	return nil
}

func deleteMedia(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	media := mediaFromResource(d)
	if err := wc.Media.Delete(context.Background(), media); err != nil {
		return fmt.Errorf("couldn't delete media: %s", err)
	}

	return nil
}

// Private helpers

func applyMediaFieldsToResource(m *wistia.Media, d *schema.ResourceData) {
	d.SetId(m.HashedId)
	d.Set("media_id", m.Id)
	d.Set("name", m.Name)
	d.Set("type", m.Type)
	d.Set("section", m.Section)
	//d.Set("thumbnail", m.Thumbnail)
	d.Set("duration", m.Duration)
	d.Set("created", m.Created)
	d.Set("updated", m.Updated)
	//d.Set("assets", m.Assets)
	//d.Set("embed_code", m.EmbedCode)
	d.Set("description", m.Description)
	d.Set("hashed_id", m.HashedId)
}

func mediaFromResource(d *schema.ResourceData) *wistia.Media {
	//thumbMap := d.Get("thumbnail").(map[string]interface{})
	//thumbnail := wistia.Thumbnail{}
	//if url, ok := thumbMap["url"]; ok {
	//	thumbnail.URL = url.(string)
	//}
	//if width, ok := thumbMap["width"]; ok {
	//	thumbnail.Width = width.(int)
	//}
	//if height, ok := thumbMap["height"]; ok {
	//	thumbnail.Height = height.(int)
	//}

	return &wistia.Media{
		Id:      d.Get("media_id").(int),
		Name:    d.Get("name").(string),
		Project: wistia.Project{HashedId: d.Get("project_id").(string)},
		Type:    d.Get("type").(string),
		Section: d.Get("section").(string),
		//Thumbnail:   thumbnail,
		Duration: d.Get("duration").(float64),
		Created:  d.Get("created").(string),
		Updated:  d.Get("updated").(string),
		//Assets:      d.Get("assets").([]wistia.Asset),
		//EmbedCode:   d.Get("embed_code").(string),
		Description: d.Get("description").(string),
		HashedId:    d.Get("hashed_id").(string),
	}
}
