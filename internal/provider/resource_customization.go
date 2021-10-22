package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wistia/terraform-provider-wistia/internal/wistia"
	"log"
)

func customizationResource() *schema.Resource {
	return &schema.Resource{
		Create:      createCustomization,
		Read:        readCustomization,
		Update:      updateCustomization,
		Delete:      deleteCustomization,
		Description: "Customize a media. See [embed options](https://wistia.com/support/developers/embed-options) for the most up-to-date documentation, including examples.",

		Schema: map[string]*schema.Schema{
			"auto_play": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the video will play as soon as it's ready.",
			},
			"controls_visible_on_load": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, controls like the big play button, playbar, volume, etc. will be visible as soon as the video is embedded. Default is true.",
			},
			"copy_link_and_thumbnail_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `If set to false, once your video is embedded on a webpage, the option to "Copy Link and Thumbnail" when you right click on your video will be removed. NOTE: If set to false, you will not be able to create a thumbnail that links to the page where the video is embedded. Default is true."`,
			},
			"do_not_track": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "By default, data for each viewing session is tracked and reported back to the Wistia servers for display in heatmaps and aggregation graphs. If you do not want to track viewing sessions, set doNotTrack to true.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Associate a specific email address with this video’s viewing sessions. This is equivalent to running video.email(email) immediately after initialization.",
			},
			"end_video_behavior": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This option determines what happens when the video ends. Possible values are default, reset, and loop.",
			},
			"fake_fullscreen": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default is false. On mobile, for certain devices (i.e. iOS), we pass the video to the native player. We do this because forcing our player to go fullscreen can cause issues with formatting, which can be a jarring experience for the viewer. This means that customizations which come with our player do not appear. You can get around this by setting this option to true.",
			},
			"fit_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is used to resize a video when there’s a discrepancy between its aspect ratio and that of its parent container. It has the effect of resizing the video independently of the Wistia player.",
			},
			"fullscreen_button": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the fullscreen button will be available as a video control",
			},
			"fullscreen_on_rotate_to_landscape": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default is true. If set to false, the video will not automatically go to true fullscreen on a mobile device. The player will rotate, and your viewer can still click on the fullscreen option after rotating.",
			},
			"google_analytics": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If you’re using Google Analytics on the page where you embed a video, the video will auto-magically send events to your Google Analytics account 📈",
			},
			"media_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The identifier of the media that's being customized.",
			},
			"muted": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the video will start in a muted state.",
			},
			"playback_rate_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to false, the playback speed controls in the settings menu will be hidden.",
			},
			"playbar": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the playbar — which includes the playhead, current time, and scrubbing functionality — will be available. If set to false, it is hidden.",
			},
			"play_button": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the big play button control will appear in the center of the video before play.",
			},
			"player_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Changes the base color of the player. Expects a hexadecimal rgb string like “ff0000” (red), “000000” (black), “ffffff” (white), or “0000ff” (blue).",
			},
			//"playlist_links": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//  Description: "The playlistLinks option lets you associate specially crafted links on your page with a video, turning them into a playlist",
			//},
			"playlist_loop": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set to true and this video has a playlist, it will loop back to the first video and replay it once the last video has finished.",
			},
			"playsinline": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set to false, your videos will play within the native mobile player instead of our own. This can be helpful if, for example, you would prefer that your mobile viewers start the video in fullscreen mode upon pressing play.",
			},
			"play_suspended_off_screen": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When a video is set to autoPlay=muted, it will pause playback when the video is out of view. For example, if the video is at the top of a page and you scroll down past it, the video will pause until you scroll back up to see the video again. To prevent a muted autoplay video from pausing when out of view, you can set this to false.",
			},
			// TODO
			//plugin[videoThumbnail][clickToPlayButton]
			"preload": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This sets the video’s preload property. Possible values are metadata, auto, none, true, and false.",
			},
			"quality_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to false, the video quality selector in the settings menu will be hidden.",
			},
			"quality_max": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting a qualityMax allows you to specify the maximum quality the video will play at. Wistia will still run bandwidth checks to test for speed, and play the highest quality version at or below the set maximum. Accepted values: 224, 360, 540, 720, 1080, 3840. Keep in mind, viewers will still be able to manually select a quality outside set maximum (using the option on the player) unless qualityControl is set to false.",
			},
			"quality_min": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting a qualityMin allows you to specify the minimum quality the video will play at. Wistia will still run bandwidth checks to test for speed, and play the highest quality version at or above the set minimum. Accepted values: 224, 360, 540, 720, 1080, 3840. Keep in mind, viewers will still be able to manually select a quality outside set minimum (using the option on the player) unless qualityControl is set to false.",
			},
			"resumable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The resumable feature causes videos to pick up where the viewer left off next time they load the page where your video is embedded. Possible values for the resumable embed option are true, false, and auto. Defaults to auto. If auto, the resumable feature will only be enabled if the video is 5 minutes or longer, is not autoplay, and is not set to loop at the end. Setting resumable to true will enable resumable regardless of those factors, and false disables resumable no matter what.",
			},
			"seo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the video’s metadata will be injected into the page’s on-page markup. Set it to false if you don’t want SEO metadata injection. For more information about how this works, check out the video SEO page. NOTE: Only the Standard and Popover embeds are capable of injecting metadata right now. The Fallback iframe embed will not inject metadata, even if seo is set to true.",
			},
			"settings_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the settings control will be available. This will allow viewers to control the quality and playback speed of the video. See qualityControl and playbackRateControl if you want control of what is available in the settings control.",
			},
			"silent_auto_play": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This option allows videos to autoplay in a muted state in contexts where normal autoplay is blocked or not supported (e.g. iOS, Safari 11+, Chrome 66+). Possible values are true, allow, and false.",
			},
			"small_play_button": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true, the small play button control will be available.",
			},
			"still_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Overrides the thumbnail image that appears before the video plays. Expects an absolute URL to an image. For best results, the image should match the exact aspect ratio of the video.",
			},
			"time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Set the time at which the video should start. Expects an integer value in seconds or string values like "5m45s". This is equivalent to running video.time(t) immediately after initialization.`,
			},
			// TODO: support thumbnailAltText
			// TODO: support foam bounds
			"video_foam": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set to true, the video will monitor the width of its parent element. When that width changes, the video will match that width and modify its height to maintain the correct aspect ratio.",
			},
			"volume": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Set the volume of the video. Expects an integer value between 0 and 1. This is equivalent to running video.volume(v) immediately after initialization. To mute the video on load, set this option to 0.",
			},
			"volume_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set to true, a volume control is available over the video. NOTE: On mobile devices where we use the native volume controls, this option has no effect.",
			},
			"wmode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set to transparent, the background behind the video player will be transparent - allowing the page color to show through - instead of black. This applies e.g. if there’s an aspect ratio discrepancy between the dimensions of the video and its container; this option is not connected to Alpha Transparency.",
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
