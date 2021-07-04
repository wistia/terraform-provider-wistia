package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/wistia/terraform-provider-wistia/wistia"
	"log"
)

func customizationResource() *schema.Resource {
	return &schema.Resource{
		Create: createCustomization,
		Read:   readCustomization,
		Update: updateCustomization,
		Delete: deleteCustomization,

		Schema: map[string]*schema.Schema{
			"auto_play": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"controls_visible_on_load": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"copy_link_and_thumbnail_enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"do_not_track": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_video_behavior": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fake_fullscreen": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fit_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fullscreen_button": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fullscreen_on_rotate_to_landscape": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"google_analytics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"media_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"muted": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"playback_rate_control": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"playbar": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"play_button": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"player_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			//"playlist_links": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			"playlist_loop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"playsinline": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"play_suspended_off_screen": {
				Type:     schema.TypeString,
				Optional: true,
			},
			//plugin[videoThumbnail][clickToPlayButton]
			"preload": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quality_control": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quality_max": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quality_min": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resumable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"seo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"settings_control": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"silent_auto_play": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"small_play_button": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"still_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO: support foam bounds
			"video_foam": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_control": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"wmode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createCustomization(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	c := customizationFromResource(d)
	c, err := wc.Customizations.Create(context.Background(), c)
	if err != nil {
		return fmt.Errorf("couldn't create Wistia customization: %s", err)
	}

	log.Printf("[TRACE] Newly created customization: %v", c)

	applyCustomizationFieldsToResource(c, d)

	return nil
}

func readCustomization(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	c, err := wc.Customizations.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("couldn't get Wistia project: %s", err)
	}

	applyCustomizationFieldsToResource(c, d)

	return nil
}

func updateCustomization(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	c := customizationFromResource(d)
	log.Printf("[TRACE] Customization before update: %v", c)
	c, err := wc.Customizations.Update(context.Background(), c)
	if err != nil {
		return fmt.Errorf("couldn't update Wistia customization: %s", err)
	}

	log.Printf("[TRACE] Customization after update: %v", c)

	applyCustomizationFieldsToResource(c, d)

	return nil
}

func deleteCustomization(d *schema.ResourceData, m interface{}) error {
	wc := m.(*wistia.Client)
	c := customizationFromResource(d)
	if err := wc.Customizations.Delete(context.Background(), c); err != nil {
		return fmt.Errorf("couldn't delete Wistia customization: %s", err)
	}

	return nil
}

// Private helpers

func applyCustomizationFieldsToResource(c *wistia.Customization, d *schema.ResourceData) {
	d.SetId(c.Media.HashedId)

	d.Set("auto_play", c.AutoPlay)
	d.Set("controls_visible_on_load", c.ControlsVisibleOnLoad)
	d.Set("copy_link_and_thumbnail_enabled", c.CopyLinkAndThumbnailEnabled)
	d.Set("do_not_track", c.DoNotTrack)
	d.Set("email", c.Email)
	d.Set("end_video_behavior", c.EndVideoBehavior)
	d.Set("fake_fullscreen", c.FakeFullscreen)
	d.Set("fit_strategy", c.FitStrategy)
	d.Set("fullscreen_button", c.FullscreenButton)
	d.Set("fullscreen_on_rotate_to_landscape", c.FullscreenOnRotateToLandscape)
	d.Set("google_analytics", c.GoogleAnalytics)
	d.Set("muted", c.Muted)
	d.Set("playback_rate_control", c.PlaybackRateControl)
	d.Set("playbar", c.Playbar)
	d.Set("play_button", c.PlayButton)
	d.Set("player_color", c.PlayerColor)
	d.Set("playlist_loop", c.PlaylistLoop)
	d.Set("playsinline", c.Playsinline)
	d.Set("play_suspended_off_screen", c.PlaySuspendedOffScreen)
	d.Set("preload", c.Preload)
	d.Set("quality_control", c.QualityControl)
	d.Set("quality_max", c.QualityMax)
	d.Set("quality_min", c.QualityMin)
	d.Set("resumable", c.Resumable)
	d.Set("seo", c.Seo)
	d.Set("settings_control", c.SettingsControl)
	d.Set("silent_auto_play", c.SilentAutoPlay)
	d.Set("small_play_button", c.SmallPlayButton)
	d.Set("still_url", c.StillUrl)
	d.Set("time", c.Time)
	d.Set("video_foam", c.VideoFoam)
	d.Set("volume", c.Volume)
	d.Set("volume_control", c.VolumeControl)
	d.Set("wmode", c.Wmode)
}

func toPointer(s string) *string {
	if s == "" {
		return nil
	}
	
	return &s
}

func customizationFromResource(d *schema.ResourceData) *wistia.Customization {
	return &wistia.Customization{
		AutoPlay:                      toPointer(d.Get("auto_play").(string)),
		ControlsVisibleOnLoad:         toPointer(d.Get("controls_visible_on_load").(string)),
		CopyLinkAndThumbnailEnabled:   toPointer(d.Get("copy_link_and_thumbnail_enabled").(string)),
		DoNotTrack:                    toPointer(d.Get("do_not_track").(string)),
		Email:                         toPointer(d.Get("email").(string)),
		EndVideoBehavior:              toPointer(d.Get("end_video_behavior").(string)),
		FakeFullscreen:                toPointer(d.Get("fake_fullscreen").(string)),
		FitStrategy:                   toPointer(d.Get("fit_strategy").(string)),
		FullscreenButton:              toPointer(d.Get("fullscreen_button").(string)),
		FullscreenOnRotateToLandscape: toPointer(d.Get("fullscreen_on_rotate_to_landscape").(string)),
		GoogleAnalytics:               toPointer(d.Get("google_analytics").(string)),
		Media:                         wistia.Media{HashedId: d.Get("media_id").(string)},
		Muted:                         toPointer(d.Get("muted").(string)),
		PlaybackRateControl:           toPointer(d.Get("playback_rate_control").(string)),
		Playbar:                       toPointer(d.Get("playbar").(string)),
		PlayButton:                    toPointer(d.Get("play_button").(string)),
		PlayerColor:                   toPointer(d.Get("player_color").(string)),
		PlaylistLoop:                  toPointer(d.Get("playlist_loop").(string)),
		Playsinline:                   toPointer(d.Get("playsinline").(string)),
		PlaySuspendedOffScreen:        toPointer(d.Get("play_suspended_off_screen").(string)),
		Preload:                       toPointer(d.Get("preload").(string)),
		QualityControl:                toPointer(d.Get("quality_control").(string)),
		QualityMax:                    toPointer(d.Get("quality_max").(string)),
		QualityMin:                    toPointer(d.Get("quality_min").(string)),
		Resumable:                     toPointer(d.Get("resumable").(string)),
		Seo:                           toPointer(d.Get("seo").(string)),
		SettingsControl:               toPointer(d.Get("settings_control").(string)),
		SilentAutoPlay:                toPointer(d.Get("silent_auto_play").(string)),
		SmallPlayButton:               toPointer(d.Get("small_play_button").(string)),
		StillUrl:                      toPointer(d.Get("still_url").(string)),
		Time:                          toPointer(d.Get("time").(string)),
		VideoFoam:                     toPointer(d.Get("video_foam").(string)),
		Volume:                        toPointer(d.Get("volume").(string)),
		VolumeControl:                 toPointer(d.Get("volume_control").(string)),
		Wmode:                         toPointer(d.Get("wmode").(string)),
	}
}
