package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Storage interface {
	Put(filename string, data []byte) error
	Get(filename string) ([]byte, error)
	Delete(filename string) error
	Exists(filename string) (bool, error)
	GetMeta(filename string) (map[string]string, error)
	GetURL(filename string) (string, error)
	GetPath(filename string) (string, error)
	Log(filename string, data []byte) error
}

type localStorage struct{}

var storageAppRoot string = "storage/app"
var storageLogRoot string = "storage/logs"
var storageCacheRoot string = "storage/cache"

func StorageInit() error {
	err := os.MkdirAll(storageAppRoot, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(storageLogRoot, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(storageCacheRoot, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

var logFile *os.File

func LogFile(filename string) *os.File {
	logFile, err := os.OpenFile(storageLogRoot+"/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	return logFile
}
func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}

func NewStorage() Storage {
	return &localStorage{}
}

func (s *localStorage) Put(filename string, data []byte) error {

	if len(data) == 0 {
		return nil // Nothing to create or save
	}
	err := os.MkdirAll(storageAppRoot, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := storageAppRoot + "/" + filename
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Create and open the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewBuffer(data))
	return err
}

func (s *localStorage) Log(filename string, data []byte) error {

	file, err := os.Create(storageLogRoot + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewBuffer(data))
	return err
}

func (s *localStorage) Get(filename string) ([]byte, error) {
	file, err := os.Open(storageAppRoot + "/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	return data, err
}

func (s *localStorage) Delete(filename string) error {
	return os.Remove(storageAppRoot + "/" + filename)
}

func (s *localStorage) Exists(filename string) (bool, error) {
	_, err := os.Stat(storageAppRoot + "/" + filename)
	if filename == "" {
		return false, fmt.Errorf("path can't be empty")
	}

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *localStorage) GetMeta(filename string) (map[string]string, error) {
	file, err := os.Open(storageAppRoot + "/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	meta := map[string]string{
		"filename":     filename,
		"size":         fmt.Sprintf("%d", info.Size()),
		"mimeType":     http.DetectContentType([]byte{}),
		"lastModified": info.ModTime().Format(time.RFC3339),
	}
	return meta, nil
}

func (s *localStorage) GetURL(filename string) (string, error) {
	url, err := url.Parse("/" + storageAppRoot + "/" + filename)
	return url.String(), err
}

func (s *localStorage) GetPath(filename string) (string, error) {
	return filepath.Join(storageAppRoot+"/", filename), nil
}
