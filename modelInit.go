package main

import (
	"archive/zip"
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func VecBaseInicialisation() map[string][]float64 {
	var archiveName string = "model.zip"
	var archiveDir string = "ArchUnzip"

	_, err := os.Stat(archiveName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(err)
			log.Println("Model archive download started")
			downloadZip(archiveName)

			if err := unzip(archiveName, archiveDir); err != nil {
				log.Println(err)
			} else {
				log.Println("Archive unziped in: ", archiveDir)
			}
		}
	} else {
		log.Println("Model archive exists:", archiveName, "; Unziped model placed in:", archiveDir)
	}

	return parseModelData(archiveDir + "\\model.txt")
}

func parseModelData(pathToModel string) map[string][]float64 {
	vecHolder := make(map[string][]float64)
	var lastWord string
	log.Println("Parse started")
	file, err := os.Open(pathToModel)
	if err != nil {
		file.Close()
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if resFloat, err := strconv.ParseFloat(scanner.Text(), 64); err == nil {
			vecHolder[lastWord] = append(vecHolder[lastWord], resFloat)
		} else {
			lastWord = scanner.Text()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal()
	}
	file.Close()

	return vecHolder
}

func downloadZip(name string) {
	var URL string = "http://vectors.nlpl.eu/repository/20/220.zip"
	response, err := http.Get(URL)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	log.Printf("err: %s", err)
}

func unzip(source string, destination string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		return err
	}

	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		filePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}
