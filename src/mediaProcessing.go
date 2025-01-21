package minefin

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type mediaFileProcessor struct {
	thumbnail_dir    string
	supportedFormats []string
}

func (m mediaFileProcessor) ScanDirectory(path string) []MediaList {
	var media_list []MediaList

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory: ", err)
	}

	for _, file := range files {
		if m.isSupported(file.Name()) {
			media_list = append(media_list, m.processFile(path, file.Name()))
		}
	}

	return media_list
}

func (m mediaFileProcessor) isSupported(filename string) bool {
	m.supportedFormats = append(m.supportedFormats, ".mp4")
	m.supportedFormats = append(m.supportedFormats, ".webm")
	m.supportedFormats = append(m.supportedFormats, ".mkv")

	for _, extension := range m.supportedFormats {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}

	return false
}

func (m mediaFileProcessor) processFile(directory string, filename string) MediaList {
	file_path := directory + "/" + filename
	file_stats, err := os.Stat(file_path)
	if err != nil {
		log.Fatal(err)
	}
	checksum := m.generateChecksum(file_path, 8192)

	fileBytes, err := ioutil.ReadFile(file_path)
	if err != nil {
		log.Fatal(err)
	}

	mime_type := http.DetectContentType(fileBytes)

	mediaList := MediaList{
		id:          checksum[:12],
		file_path:   file_path,
		file_name:   filename,
		file_size:   int(file_stats.Size()),
		mime_type:   mime_type,
		created_at:  file_stats.ModTime(),
		modified_at: file_stats.ModTime(),
		checksum:    checksum,
	}

	return mediaList
}

func (m mediaFileProcessor) generateChecksum(path string, chunkSize int) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()

	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return ""
		}
		if n == 0 {
			break
		}

		hash.Write(buffer[:n])
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
