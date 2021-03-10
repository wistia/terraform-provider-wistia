package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/wistia/terraform-provider-wistia/wistia"
	"net/http"
	"time"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: configureProvider,
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WISTIA_ACCESS_TOKEN", nil),
				Description: "Wistia access token with read, update, delete, and upload permissions",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WISTIA_ENV", "production"),
				Description: "Wistia environment to use [production (default), staging]",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wistia_media":   mediaResource(),
			"wistia_project": projectResource(),
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	accessToken := d.Get("access_token").(string)
	environment := d.Get("environment").(string)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	wistiaClient := wistia.NewClient(httpClient, accessToken)
	if environment == "staging" {
		wistiaClient.APIBaseEndpoint = "https://api.wistia.st/v1/"
		wistiaClient.UploadBaseEndpoint = "https://upload-v2.wistia.st/"
	}
	return wistiaClient, nil
}
