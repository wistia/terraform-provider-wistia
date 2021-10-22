package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Description: "A Wistia media. See the [API documentation](https://wistia.com/support/developers/data-api#medias) for more details.",

		Schema: map[string]*schema.Schema{
			"file": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file", "url"},
				Description:  "A path to a file on disk that will be uploaded to Wistia.",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"file", "url"},
				Description:  "A URL to a file that will be uploaded to Wistia.",
			},
			"media_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique numeric identifier for the media within the system.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the media.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The identifier for the Wistia project that will host this media.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A string representing what type of media this is. Values can be Video, Audio, Image, PdfDocument, MicrosoftOfficeDocument, Swf, or UnknownType.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Post upload processing status. There are four statuses: queued, processing, ready, and failed.",
			},
			"section": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The title of the section in which the media appears. This attribute is omitted if the media is not in a section (default).",
			},
			// TODO
			//"thumbnail": {
			//	Type:        schema.TypeMap,
			//	Computed:    true,
			//  Description: "An object representing the thumbnail for this media. The attributes are URL, width, and height.
			//},
			"duration": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Specifies the length (in seconds) for audio and video files. Specifies the number of pages in the document. Omitted for other types of media.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the media was originally uploaded.",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the media was last changed.",
			},
			// TODO
			//"assets": {
			//	Type:        schema.TypeList,
			//	Elem:        wistia.Asset{},
			//	Computed:    true,
			//  Description: "An array of the assets available for this media.",
			//},
			//"embed_code": {
			//	Type:       schema.TypeString,
			//	Computed:   true,
			//	Deprecated: "If you want to programmatically embed videos, follow the \"construct an embed code\" guide.",
			//},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description for the media which usually appears near teh top of the sidebar on the media's page.",
			},
			"hashed_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique alphanumeric identifier for this media.",
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
