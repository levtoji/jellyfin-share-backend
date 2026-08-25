package proxy

import (
	"net/url"
	"strings"
	"testing"

	"github.com/jellyfin-share/jellyfin-share-backend/internal/config"
	"github.com/jellyfin-share/jellyfin-share-backend/internal/jellyfin"
)

func newTestStreamProxy(baseURL, apiKey string) *StreamProxy {
	jf := jellyfin.NewClient(baseURL, apiKey)
	cfg := &config.Config{}
	return NewStreamProxy(nil, jf, cfg)
}

func TestBuildJellyfinStreamURL_MasterM3u8(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "master.m3u8", "")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()

	if !strings.HasSuffix(parsed.Path, "/Videos/item-123/master.m3u8") {
		t.Errorf("unexpected path: %s", parsed.Path)
	}
	if params.Get("api_key") != "test-api-key" {
		t.Errorf("expected api_key=test-api-key, got %s", params.Get("api_key"))
	}
	if params.Get("MediaSourceId") != "item-123" {
		t.Errorf("expected MediaSourceId=item-123, got %s", params.Get("MediaSourceId"))
	}
	if params.Get("DeviceId") != "jfshare-backend" {
		t.Errorf("expected DeviceId=jfshare-backend, got %s", params.Get("DeviceId"))
	}
	if params.Get("VideoCodec") != "h264" {
		t.Errorf("expected VideoCodec=h264, got %s", params.Get("VideoCodec"))
	}
	if params.Get("AllowVideoStreamCopy") != "false" {
		t.Errorf("expected AllowVideoStreamCopy=false, got %s", params.Get("AllowVideoStreamCopy"))
	}
	if params.Get("AudioCodec") != "aac" {
		t.Errorf("expected AudioCodec=aac, got %s", params.Get("AudioCodec"))
	}
}

func TestBuildJellyfinStreamURL_MasterM3u8_WithAudioLanguage(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "master.m3u8", "AudioLanguage=eng")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()

	if params.Get("AudioLanguage") != "eng" {
		t.Errorf("expected AudioLanguage=eng, got %s", params.Get("AudioLanguage"))
	}
	if params.Get("VideoCodec") != "h264" {
		t.Errorf("expected VideoCodec=h264, got %s", params.Get("VideoCodec"))
	}
}

func TestBuildJellyfinStreamURL_MasterM3u8_WithItemId(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "master.m3u8", "itemId=episode-456&AudioLanguage=eng")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()

	if params.Get("itemId") != "episode-456" {
		t.Errorf("expected itemId=episode-456, got %s", params.Get("itemId"))
	}
	if params.Get("AudioLanguage") != "eng" {
		t.Errorf("expected AudioLanguage=eng, got %s", params.Get("AudioLanguage"))
	}
	if params.Get("MediaSourceId") != "item-123" {
		t.Errorf("expected MediaSourceId=item-123, got %s", params.Get("MediaSourceId"))
	}
}

func TestBuildJellyfinStreamURL_SubPlaylist(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "hls1/main/index.m3u8", "SomeParam=value")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	if !strings.HasSuffix(parsed.Path, "/Videos/item-123/hls1/main/index.m3u8") {
		t.Errorf("unexpected path: %s", parsed.Path)
	}

	params := parsed.Query()
	if params.Get("SomeParam") != "value" {
		t.Errorf("expected SomeParam=value, got %s", params.Get("SomeParam"))
	}
}

func TestBuildJellyfinStreamURL_TsSegment_RemovesAudioCodec(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "hls1/main/42.ts", "AudioCodec=m3u8&DeviceId=jfshare-backend")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()

	if params.Get("AudioCodec") != "" {
		t.Errorf("expected AudioCodec to be removed for .ts segments, got %s", params.Get("AudioCodec"))
	}
	if params.Get("DeviceId") != "jfshare-backend" {
		t.Errorf("expected DeviceId=jfshare-backend, got %s", params.Get("DeviceId"))
	}
}

func TestBuildJellyfinStreamURL_M4sSegment_RemovesAudioCodec(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "hls1/main/42.m4s", "AudioCodec=aac")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()
	if params.Get("AudioCodec") != "" {
		t.Errorf("expected AudioCodec to be removed for .m4s segments, got %s", params.Get("AudioCodec"))
	}
}

func TestBuildJellyfinStreamURL_GenericStream(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "stream", "")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	if !strings.HasSuffix(parsed.Path, "/Videos/item-123/stream") {
		t.Errorf("unexpected path: %s", parsed.Path)
	}

	params := parsed.Query()
	if params.Get("Static") != "true" {
		t.Errorf("expected Static=true, got %s", params.Get("Static"))
	}
	if params.Get("mediaSourceId") != "item-123" {
		t.Errorf("expected mediaSourceId=item-123, got %s", params.Get("mediaSourceId"))
	}
}

func TestBuildJellyfinStreamURL_NoDuplicateApiKey(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	result := proxy.buildJellyfinStreamURL("item-123", "master.m3u8", "api_key=already-set")

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()
	if params.Get("api_key") != "already-set" {
		t.Errorf("expected api_key=already-set (not overwritten), got %s", params.Get("api_key"))
	}
}

func TestBuildJellyfinStreamURL_MultipleParams(t *testing.T) {
	proxy := newTestStreamProxy("http://jellyfin:8096", "test-api-key")

	query := "VideoCodec=h264&AllowVideoStreamCopy=false&VideoBitrate=4000000&MaxWidth=1920&MaxHeight=1080&EnableTonemapping=true&AudioLanguage=eng"
	result := proxy.buildJellyfinStreamURL("item-123", "master.m3u8", query)

	parsed, err := url.Parse(result)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	params := parsed.Query()

	expectedParams := map[string]string{
		"VideoCodec":           "h264",
		"AllowVideoStreamCopy": "false",
		"VideoBitrate":         "4000000",
		"MaxWidth":             "1920",
		"MaxHeight":            "1080",
		"EnableTonemapping":    "true",
		"AudioLanguage":        "eng",
	}

	for key, expectedVal := range expectedParams {
		if params.Get(key) != expectedVal {
			t.Errorf("expected %s=%s, got %s", key, expectedVal, params.Get(key))
		}
	}
}
