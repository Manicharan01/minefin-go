package main

import (
	"fmt"
	"minefin/minefin"
)

func main() {
	mediaProcessor := minefin.MediaFileProcessor{
		Thumbnail_dir:    "/home/charan/Downloads/thumbnails",
		SupportedFormats: []string{"avi", "mp4", "webm", "mkv"},
		Directory:        "/LocalDisk/Movies",
	}

	fmt.Println("Scanning the given direcotory")
	// mediaProcessor.DirWathcer("/home/charan/Videos/YouTube")
	// mediaProcessor.PostgreSQL()
	mediaProcessor.ScanDirectory()
	// mediaProcessor.PostgreSQL()
}
