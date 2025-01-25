package minefin

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func (m MediaFileProcessor) DirWathcer(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error while creating a dir watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			// provide the list of events to monitor
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("File created:", event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("File removed:", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("File renamed:", event.Name)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("File permissions modified:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatalf("Error while monitoring the dir: %v", err)
	}
	<-done
}
