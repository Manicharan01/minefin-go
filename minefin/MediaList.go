package minefin

import (
	"time"
)

type MediaList struct {
	id          string
	file_path   string
	file_name   string
	file_size   int
	mime_type   string
	created_at  time.Time
	modified_at time.Time
	// checksum    string

	//Media metadata
	duration       int
	width          int
	height         int
	bitrate        int
	codec          string
	frame_rate     float64
	audio_channels int
	audio_codec    string

	//Additional Metadata
	title          string
	thumbnail_path string
}
