package minefin

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"mime"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/u2takey/ffmpeg-go"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type MediaFileProcessor struct {
	Thumbnail_dir    string
	SupportedFormats []string
	Directory        string
}

func (m MediaFileProcessor) ScanDirectory() {
	libraryManager := LibraryManager{
		LibraryPath: m.Directory,
	}

	tableNames := libraryManager.GetTableNames()
	fmt.Println(tableNames)
	if !contains(tableNames, "users") {
		libraryManager.CreateTable("users")
	}
	if !contains(tableNames, "mediaitems") {
		libraryManager.CreateTable("mediaitems")
	}
	if !contains(tableNames, "usermediaprogress") {
		libraryManager.CreateTable("usermediaprogress")
	}
}

func (m MediaFileProcessor) getMetadata(path string) []MediaList {
	var media_list []MediaList

	fmt.Println("Reading the directory")
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory: ", err)
	}

	for _, file := range files {
		count := 1
		fmt.Printf("Checking and Processing file number %d", count)
		if m.isSupported(file.Name()) {
			media_list = append(media_list, m.processFile(path, file.Name()))
		}
		fmt.Println("File has been added to MediaList")

		count += 1
	}

	return media_list
}

func (m MediaFileProcessor) isSupported(filename string) bool {
	fmt.Println("Checking if the file is Supported or not")
	for _, extension := range m.SupportedFormats {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}

	return false
}

func (m MediaFileProcessor) processFile(directory string, filename string) MediaList {
	fmt.Println("Processing and adding the metadata for the given file")
	file_path := directory + "/" + filename
	fmt.Println(file_path)
	fmt.Println("Getting File stats")
	file_stats, err := os.Stat(file_path)
	if err != nil {
		log.Fatal(err)
	}
	checksum := m.generateChecksum(file_path, 1073741824)

	fmt.Println("Getting MIME types")
	extension := strings.Split(file_path, ".")
	mime_type := mime.TypeByExtension(extension[len(extension)-1])
	id_int := rand.Int()
	id := strconv.Itoa(id_int)

	mediaList := MediaList{
		id:          id,
		file_path:   file_path,
		file_name:   filename,
		file_size:   int(file_stats.Size()),
		mime_type:   mime_type,
		created_at:  file_stats.ModTime(),
		modified_at: file_stats.ModTime(),
		checksum:    checksum,
	}

	fmt.Println("Probing the given file")
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	probeData, err := ffprobe.ProbeURL(ctx, file_path)
	if err != nil {
		log.Panicf("Error getting the probeData: %v", err)
	}
	fmt.Println("Getting metadata of Audio and Video streams")

	video_stream := probeData.FirstVideoStream()
	audio_stream := probeData.FirstAudioStream()

	if video_stream != nil {
		fmt.Println("Video stream metadata")
		mediaList.width = video_stream.Width
		mediaList.height = video_stream.Height
		mediaList.codec = video_stream.CodecName
		framerate, _ := strconv.ParseFloat(video_stream.RFrameRate, 8)
		mediaList.frame_rate = framerate
	}

	if audio_stream != nil {
		fmt.Println("Audio stream metadata")
		mediaList.audio_channels = audio_stream.Channels
		mediaList.audio_codec = audio_stream.CodecName
	}

	formats := probeData.Format
	mediaList.duration = int(formats.Duration())
	bitRate, _ := strconv.Atoi(formats.BitRate)
	mediaList.bitrate = bitRate

	mediaList.title, _ = formats.TagList.GetString("title")

	if video_stream != nil || strings.HasPrefix(mime_type, "image/") {
		mediaList.thumbnail_path = m.generateThumbnail(file_path, mediaList.id)
	}
	fmt.Println("Metadata is generated")

	return mediaList
}

func (m MediaFileProcessor) generateChecksum(path string, chunkSize int) string {
	fmt.Println("Generating checksum")
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
	fmt.Println("Checksum is generated")

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (m MediaFileProcessor) generateThumbnail(file_path string, file_id string) string {
	fmt.Println("Generating Thumbnail")
	thumbnailPath := m.Thumbnail_dir + "/" + file_id + ".jpg"
	var image_formats []string
	var isImage bool
	image_formats = append(image_formats, ".jpg")
	image_formats = append(image_formats, ".jpeg")
	image_formats = append(image_formats, ".png")
	image_formats = append(image_formats, ".gif")

	for _, format := range image_formats {
		if strings.HasSuffix(file_path, format) {
			isImage = true
			break
		}
	}

	if isImage {
		file, err := os.Open(file_path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		thumbnail := resize.Thumbnail(320, 180, img, resize.Lanczos3)

		outFile, err := os.Create(thumbnailPath)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		err = jpeg.Encode(outFile, thumbnail, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Generating thumbnail using video_stream")

		reader := ExampleReadFrameAsJpeg(file_path, 1)
		img, _, err := image.Decode(reader)
		if err != nil {
			log.Fatalf("Error generating thumbnail: %v", err)
		}

		outFile, err := os.Create(thumbnailPath)
		if err != nil {
			log.Fatalf("Error generating thumbnail: %v", err)
		}
		defer outFile.Close()

		err = jpeg.Encode(outFile, img, nil)
		if err != nil {
			log.Fatalf("Error generating thumbnail: %v", err)
		}
	}
	fmt.Println("Generated thumbnail and added to thumbnail_dir")

	return thumbnailPath
}

func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(inFileName).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 10, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}
