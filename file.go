package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func saveStatus(fileName string, from, to int) {
	baseFileName := filepath.Base(fileName)
	f, err := os.Create(filepath.Join(getHomeDirectoryPath(runtime.GOOS), "ltbr", "progress", baseFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s|%d|%d", fileName, from, to))
}

func getFileNameFromLatest(filePath string) (LatestFile, error) {
	baseFileName := filepath.Base(filePath)
	latestFilePath := filepath.Join(getHomeDirectoryPath(runtime.GOOS), "ltbr", "progress", baseFileName)
	latestFile := LatestFile{FileName: filePath, From: 0, To: Advance}
	if !exists(latestFilePath) {
		return latestFile, nil
	}

	f, err := os.Open(latestFilePath)
	if err != nil {
		return LatestFile{}, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	latestFileFields := strings.Split(string(content), "|")
	if len(latestFileFields) != DBFileRequiredNumbermields {
		return LatestFile{}, fmt.Errorf("wrong format in '%s'", latestFilePath)
	}

	latestFile.FileName = latestFileFields[0]
	fromInt, _ := strconv.ParseInt(latestFileFields[1], 10, 32)
	toInt, _ := strconv.ParseInt(latestFileFields[2], 10, 32)
	latestFile.From = int(fromInt)
	latestFile.To = int(toInt)

	return latestFile, nil
}

func dirExists(dirPath string) bool {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func readLines(file io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

func appendLineToFile(filePath, line, sep string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n%s", sep, line)); err != nil {
		panic(err)
	}
}

func createDirectories() error {
	ltbrDir := filepath.Join(getHomeDirectoryPath(runtime.GOOS), "ltbr")
	if err := createDirectory(ltbrDir); err != nil {
		return err
	}
	return createDir("notes", "quotes", "progress")
}

func createDir(dirs ...string) error {
	for _, dir := range dirs {
		if err := createDirectory(filepath.Join(getHomeDirectoryPath(runtime.GOOS), "ltbr", dir)); err != nil {
			return err
		}
	}
	return nil
}

func getHomeDirectoryPath(opSystem string) string {
	if opSystem == "windows" {
		return os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

func getDirectoryNameForFile(dirType, fileName string) string {
	absoluteFilePath, _ := filepath.Abs(fileName)
	baseFileName := path.Base(absoluteFilePath)

	baseFileName = sanitizeFileName(baseFileName)
	notesDir := filepath.Join(getHomeDirectoryPath(runtime.GOOS), "ltbr", dirType, baseFileName)

	return notesDir
}

func createDirectory(dirPath string) error {
	if !dirExists(dirPath) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}
