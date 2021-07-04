package wistia

import (
	"context"
	"fmt"
	"net/http"
)

type CustomizationsProvider provider

type Customization struct {
	AutoPlay                      *string `json:"autoPlay"`
	ControlsVisibleOnLoad         *string `json:"controlsVisibleOnLoad"`
	CopyLinkAndThumbnailEnabled   *string `json:"copyLinkAndThumbnailEnabled"`
	DoNotTrack                    *string `json:"doNotTrack"`
	Email                         *string `json:"email"`
	EndVideoBehavior              *string `json:"endVideoBehavior"`
	FakeFullscreen                *string `json:"fakeFullscreen"`
	FitStrategy                   *string `json:"fitStrategy"`
	FullscreenButton              *string `json:"fullscreenButton"`
	FullscreenOnRotateToLandscape *string `json:"fullscreenOnRotateToLandscape"`
	GoogleAnalytics               *string `json:"googleAnalytics"`
	Media                         Media  `json:"-"`
	Muted                         *string `json:"muted"`
	PlaybackRateControl           *string `json:"playbackRateControl"`
	Playbar                       *string `json:"playbar"`
	PlayButton                    *string `json:"playButton"`
	PlayerColor                   *string `json:"playerColor"`
	PlaylistLoop                  *string `json:"playlistLoop"`
	Playsinline                   *string `json:"playsinline"`
	PlaySuspendedOffScreen        *string `json:"playSuspendedOffScreen"`
	Preload                       *string `json:"preload"`
	QualityControl                *string `json:"qualityControl"`
	QualityMax                    *string `json:"qualityMax"`
	QualityMin                    *string `json:"qualityMin"`
	Resumable                     *string `json:"resumable"`
	Seo                           *string `json:"seo"`
	SettingsControl               *string `json:"settingsControl"`
	SilentAutoPlay                *string `json:"silentAutoPlay"`
	SmallPlayButton               *string `json:"smallPlayButton"`
	StillUrl                      *string `json:"stillUrl"`
	Time                          *string `json:"time"`
	VideoFoam                     *string `json:"videoFoam"`
	Volume                        *string `json:"volume"`
	VolumeControl                 *string `json:"volumeControl"`
	Wmode                         *string `json:"wmode"`
}

func (cp *CustomizationsProvider) Create(ctx context.Context, c *Customization) (*Customization, error) {
	createdCustomization := &Customization{Media: c.Media}
	url := cp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s/customizations.json", c.Media.HashedId)
	_, err := cp.client.request(ctx, http.MethodPost, url, c, createdCustomization)
	if err != nil {
		return nil, err
	}
	return createdCustomization, nil
}

func (cp *CustomizationsProvider) Get(ctx context.Context, mediaId string) (*Customization, error) {
	customization := &Customization{Media: Media{HashedId: mediaId}}
	url := cp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s/customizations.json", mediaId)
	_, err := cp.client.request(ctx, http.MethodGet, url, nil, customization)
	if err != nil {
		return nil, err
	}
	return customization, nil
}

func (cp *CustomizationsProvider) Update(ctx context.Context, c *Customization) (*Customization, error) {
	updatedCustomization := &Customization{Media: c.Media}
	url := cp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s/customizations.json", c.Media.HashedId)
	_, err := cp.client.request(ctx, http.MethodPut, url, c, updatedCustomization)
	if err != nil {
		return nil, err
	}
	return updatedCustomization, nil
}

func (cp *CustomizationsProvider) Delete(ctx context.Context, c *Customization) error {
	url := cp.client.APIBaseEndpoint + fmt.Sprintf("medias/%s/customizations.json", c.Media.HashedId)
	if _, err := cp.client.request(ctx, http.MethodDelete, url, nil, nil); err != nil {
		return err
	}

	return nil
}
