package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wistia/terraform-provider-wistia/internal/wistia"
	"net/http"
	"time"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New() func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider{
			ConfigureContextFunc: configureProvider,
			ResourcesMap: map[string]*schema.Resource{
				"wistia_media":               mediaResource(),
				"wistia_media_customization": customizationResource(),
				"wistia_project":             projectResource(),
			},
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
		}
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	accessToken := d.Get("access_token").(string)
	environment := d.Get("environment").(string)
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	wistiaClient := wistia.NewClient(httpClient, accessToken)
	if environment == "staging" {
		wistiaClient.APIBaseEndpoint = "https://api.wistia.st/v1/"
		wistiaClient.UploadBaseEndpoint = "https://upload-v2.wistia.st/"
	}
	return wistiaClient, nil
}
