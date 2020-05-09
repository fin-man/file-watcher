package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fin-man/file-watcher/watcher"

	"github.com/fin-man/finance-manager/clients/recordcreator"
	"github.com/fin-man/finance-manager/csvprocessors"
	"github.com/fin-man/finance-manager/filemanager"

	"github.com/fin-man/finance-manager/categories"
)

func main() {
	log.Println("Starting a new filewatcher ")
	fw := watcher.NewFileWatcher()

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	fullPath := pwd + "/data/transactions"

	fw.Watch(fullPath, ProcessFile)
}

func ProcessFile(data ...interface{}) error {
	fmt.Printf("FilePath : %s \n", data[0])
	fmt.Printf("FileName : %s \n", data[1])
	// fileName := data[1].(string)
	filePath := data[0].(string)

	recordCreator := recordcreator.NewRecordCreator()

	err := HandleOverall(filePath, recordCreator)
	if err != nil {
		return err
	}

	log.Printf("Unknown File Found ..")

	return nil
}

func HandleOverall(filePath string, recordCreator *recordcreator.RecordCreator) error {
	fm := filemanager.FileManager{}
	file, err := fm.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}

	defer file.Close()

	records := []*categories.NormalizedTransaction{}

	csvProcessor := csvprocessors.NewCSVprocessor()

	err = csvProcessor.Unmarshal(file, &records)

	fmt.Println(records)
	if err != nil {

		//file is prolly dont match the format
		return err
	}

	for _, v := range records {

		_, ok := categories.OverallTransactionTypes[string(v.Category)]
		if !ok {
			return fmt.Errorf("Invalid Category for record : %v", v)
		}
		err = recordCreator.CreateNewRecord(v)
		if err != nil {
			log.Println(err)
		}
	}

	return nil

}
