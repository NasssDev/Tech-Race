package models

import "time"

// Cloudinary Types

type AssetData struct {
	StatusCode int64           `json:"statusCode"`
	Message    string          `json:"message"`
	Data       InstructionData `json:"data"`
}

type InstructionData struct {
	Data DataData `json:"data"`
}

type DataData struct {
	AssetID               string      `json:"asset_id"`
	PublicID              string      `json:"public_id"`
	AssetFolder           string      `json:"asset_folder"`
	DisplayName           string      `json:"display_name"`
	Version               int64       `json:"version"`
	VersionID             string      `json:"version_id"`
	Signature             string      `json:"signature"`
	Width                 int64       `json:"width"`
	Height                int64       `json:"height"`
	Format                string      `json:"format"`
	ResourceType          string      `json:"resource_type"`
	CreatedAt             time.Time   `json:"created_at"`
	Bytes                 int64       `json:"bytes"`
	Type                  string      `json:"type"`
	Etag                  string      `json:"etag"`
	URL                   string      `json:"url"`
	SecureURL             string      `json:"secure_url"`
	PlaybackURL           string      `json:"playback_url"`
	AccessMode            string      `json:"access_mode"`
	Overwritten           bool        `json:"overwritten"`
	OriginalFilename      string      `json:"original_filename"`
	Eager                 interface{} `json:"eager"`
	ResponsiveBreakpoints interface{} `json:"responsive_breakpoints"`
	HookExecution         interface{} `json:"hook_execution"`
	Error                 Error       `json:"error"`
	Response              Response    `json:"Response"`
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	APIKey           string        `json:"api_key"`
	AssetID          string        `json:"asset_id"`
	Audio            Audio         `json:"audio"`
	BitRate          int64         `json:"bit_rate"`
	Bytes            int64         `json:"bytes"`
	CreatedAt        time.Time     `json:"created_at"`
	Duration         float64       `json:"duration"`
	Etag             string        `json:"etag"`
	Folder           string        `json:"folder"`
	Format           string        `json:"format"`
	FrameRate        int64         `json:"frame_rate"`
	Height           int64         `json:"height"`
	NbFrames         int64         `json:"nb_frames"`
	OriginalFilename string        `json:"original_filename"`
	Overwritten      bool          `json:"overwritten"`
	Pages            int64         `json:"pages"`
	Placeholder      bool          `json:"placeholder"`
	PlaybackURL      string        `json:"playback_url"`
	PublicID         string        `json:"public_id"`
	ResourceType     string        `json:"resource_type"`
	Rotation         int64         `json:"rotation"`
	SecureURL        string        `json:"secure_url"`
	Signature        string        `json:"signature"`
	Tags             []interface{} `json:"tags"`
	Type             string        `json:"type"`
	URL              string        `json:"url"`
	Version          int64         `json:"version"`
	VersionID        string        `json:"version_id"`
	Video            VideoData     `json:"video"`
	Width            int64         `json:"width"`
}

type Audio struct {
}

type VideoData struct {
	BitRate   string `json:"bit_rate"`
	Codec     string `json:"codec"`
	Level     int64  `json:"level"`
	PixFormat string `json:"pix_format"`
	Profile   string `json:"profile"`
	TimeBase  string `json:"time_base"`
}

// -- END Cloudinary Types --
