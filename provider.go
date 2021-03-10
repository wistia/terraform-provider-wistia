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
				Description: "Wistia access token",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WISTIA_ENV", "production"),
				Description: "Wistia environment to use",
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
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	wistiaClient := wistia.NewClient(httpClient, accessToken)
	return wistiaClient, nil
}
