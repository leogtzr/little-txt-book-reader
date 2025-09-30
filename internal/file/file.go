package file

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"textreader/internal/model"
)

type ProgressEntry struct {
	FileName   string   `json:"file_name"`
	From       int      `json:"from"`
	To         int      `json:"to"`
	Vocabulary []string `json:"vocabulary"`
}

func getProgressFilePath() string {
	return filepath.Join(GetHomeDirectoryPath(runtime.GOOS), "ltbr", "progress.json")
}

func SaveStatus(fileName string, from, to int, state *model.AppState) error {
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	key := hashPath(absPath)
	progressFile := getProgressFilePath()

	data := make(map[string]ProgressEntry)
	if _, err := os.Stat(progressFile); err == nil {
		content, err := os.ReadFile(progressFile)
		if err != nil {
			return fmt.Errorf("failed to read progress file: %w", err)
		}
		if err := json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal progress: %w", err)
		}
	}

	data[key] = ProgressEntry{
		FileName:   absPath,
		From:       from,
		To:         to,
		Vocabulary: state.Vocabulary,
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}
	if err := os.WriteFile(progressFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write progress file: %w", err)
	}
	return nil
}

func GetFileNameFromLatest(filePath string, state *model.AppState) (model.LatestFile, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return model.LatestFile{}, fmt.Errorf("failed to get absolute path: %w", err)
	}
	key := hashPath(absPath)
	progressFile := getProgressFilePath()

	if _, err := os.Stat(progressFile); os.IsNotExist(err) {
		return model.LatestFile{FileName: absPath, From: 0, To: state.Advance}, nil
	}

	content, err := os.ReadFile(progressFile)
	if err != nil {
		return model.LatestFile{}, fmt.Errorf("failed to read progress file: %w", err)
	}

	data := make(map[string]ProgressEntry)
	if err := json.Unmarshal(content, &data); err != nil {
		return model.LatestFile{}, fmt.Errorf("failed to unmarshal progress: %w", err)
	}

	entry, ok := data[key]
	if !ok {
		return model.LatestFile{FileName: absPath, From: 0, To: state.Advance}, nil
	}

	state.Vocabulary = entry.Vocabulary // Load vocabulary into state
	return model.LatestFile{
		FileName: entry.FileName,
		From:     entry.From,
		To:       entry.To,
	}, nil
}

func hashPath(path string) string {
	h := md5.Sum([]byte(path))
	return hex.EncodeToString(h[:])
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

func ReadLines(file io.Reader) ([]string, error) {
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

func AppendLineToFile(filePath, line, sep string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n%s", sep, line)); err != nil {
		panic(err)
	}
}

func CreateDirectories() error {
	ltbrDir := filepath.Join(GetHomeDirectoryPath(runtime.GOOS), "ltbr")
	if err := createDirectory(ltbrDir); err != nil {
		return err
	}
	return createDir("notes", "quotes", "progress", "vocabulary")
}

func createDir(dirs ...string) error {
	for _, dir := range dirs {
		if err := createDirectory(filepath.Join(GetHomeDirectoryPath(runtime.GOOS), "ltbr", dir)); err != nil {
			return err
		}
	}
	return nil
}

func GetHomeDirectoryPath(opSystem string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

func GetDirectoryNameForFile(dirType, fileName string) string {
	absoluteFilePath, _ := filepath.Abs(fileName)
	baseFileName := path.Base(absoluteFilePath)

	baseFileName = SanitizeFileName(baseFileName)
	notesDir := filepath.Join(GetHomeDirectoryPath(runtime.GOOS), "ltbr", dirType, baseFileName)

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

func SanitizeFileName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, " ", "")
	return fileName
}
