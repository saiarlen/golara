package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Storage interface defines Laravel-style storage methods
type Storage interface {
	Put(path string, contents []byte) error
	PutFile(path string, file *multipart.FileHeader) error
	Get(path string) ([]byte, error)
	Exists(path string) bool
	Delete(path string) error
	Copy(from, to string) error
	Move(from, to string) error
	Size(path string) (int64, error)
	Files(directory string) ([]string, error)
	MakeDirectory(path string) error
	DeleteDirectory(path string) error
	URL(path string) string
}

// LocalStorage implements Laravel-style local file storage
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Put stores content at the given path
func (s *LocalStorage) Put(path string, contents []byte) error {
	fullPath := filepath.Join(s.basePath, path)
	dir := filepath.Dir(fullPath)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(fullPath, contents, 0644)
}

// PutFile stores an uploaded file
func (s *LocalStorage) PutFile(path string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	
	fullPath := filepath.Join(s.basePath, path)
	dir := filepath.Dir(fullPath)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, src)
	return err
}

// Get retrieves content from the given path
func (s *LocalStorage) Get(path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.ReadFile(fullPath)
}

// Exists checks if a file exists
func (s *LocalStorage) Exists(path string) bool {
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// Delete removes a file
func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// Copy copies a file from source to destination
func (s *LocalStorage) Copy(from, to string) error {
	fromPath := filepath.Join(s.basePath, from)
	toPath := filepath.Join(s.basePath, to)
	
	dir := filepath.Dir(toPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()
	
	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, src)
	return err
}

// Move moves a file from source to destination
func (s *LocalStorage) Move(from, to string) error {
	fromPath := filepath.Join(s.basePath, from)
	toPath := filepath.Join(s.basePath, to)
	
	dir := filepath.Dir(toPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.Rename(fromPath, toPath)
}

// Size returns the size of a file
func (s *LocalStorage) Size(path string) (int64, error) {
	fullPath := filepath.Join(s.basePath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Files returns all files in a directory
func (s *LocalStorage) Files(directory string) ([]string, error) {
	fullPath := filepath.Join(s.basePath, directory)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}
	return files, nil
}

// MakeDirectory creates a directory
func (s *LocalStorage) MakeDirectory(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.MkdirAll(fullPath, 0755)
}

// DeleteDirectory removes a directory and its contents
func (s *LocalStorage) DeleteDirectory(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.RemoveAll(fullPath)
}

// URL returns the URL for a file
func (s *LocalStorage) URL(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), strings.TrimLeft(path, "/"))
}

// StorageManager manages different storage disks (Laravel-style)
type StorageManager struct {
	disks map[string]Storage
}

// NewStorageManager creates a new storage manager
func NewStorageManager() *StorageManager {
	return &StorageManager{
		disks: make(map[string]Storage),
	}
}

// Disk returns a storage disk
func (sm *StorageManager) Disk(name string) Storage {
	if disk, exists := sm.disks[name]; exists {
		return disk
	}
	return sm.disks["local"] // Default to local
}

// AddDisk adds a storage disk
func (sm *StorageManager) AddDisk(name string, storage Storage) {
	sm.disks[name] = storage
}

// Helper functions (Laravel-style)
func Store(file *multipart.FileHeader, path string) (string, error) {
	storage := NewLocalStorage("storage/app", "/storage")
	
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(path, filename)
	
	err := storage.PutFile(fullPath, file)
	if err != nil {
		return "", err
	}
	
	return fullPath, nil
}

func StoreAs(file *multipart.FileHeader, path, name string) error {
	storage := NewLocalStorage("storage/app", "/storage")
	fullPath := filepath.Join(path, name)
	return storage.PutFile(fullPath, file)
}