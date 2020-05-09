package watcher

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	Watcher *fsnotify.Watcher
}

func NewFileWatcher() *FileWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &FileWatcher{
		Watcher: watcher,
	}
}

func (f *FileWatcher) Watch(volume string, callback func(data ...interface{}) error) {
	log.Println("Starting the watch .. ")
	fmt.Println("watching : ", volume)
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-f.Watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
					fileName := f.ExtractFileName(event.Name)
					filePath := f.ExtractFilePath(event.Name)
					err := callback(filePath, fileName) // make this get executed in a go routine with a channel
					if err != nil {
						log.Printf("ERROR : %v \n", err)
					}
				}
			case err, ok := <-f.Watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := f.Watcher.Add(volume)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func (f *FileWatcher) ExtractFileName(event string) string {
	fileEventSplit := strings.Split(event, " ")
	filePathSplit := strings.Split(fileEventSplit[len(fileEventSplit)-1], "/")
	fileName := filePathSplit[len(filePathSplit)-1]

	return fileName
}

func (f *FileWatcher) ExtractFilePath(event string) string {
	fileEventSplit := strings.Split(event, " ")

	filePath := fileEventSplit[len(fileEventSplit)-1]
	return filePath
}
