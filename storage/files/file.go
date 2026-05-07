package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"main/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	base_path string
}

func New(path string) Storage {
	return Storage{path}
}

const DefaultPermission = 0774

func (s Storage) Save(page *storage.Page) error {
	file_path := filepath.Join(s.base_path, page.UserName)
	if err := os.MkdirAll(file_path, DefaultPermission); err != nil {
		return fmt.Errorf("Can't create directory %w", err)
	}

	file_name, err := FileName(page)
	if err != nil {
		return fmt.Errorf("Can't create filename %w", err)
	}

	file_path = filepath.Join(file_path, file_name)

	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("Can't create file %w", err)
	}

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("Can't convert file to gob %w", err)
	}
	file.Close()
	return nil
}

func (s Storage) PickRandom(user_name string) (*storage.Page, error) {
	file_path := filepath.Join(s.base_path, user_name)

	files, err := os.ReadDir(file_path)
	if err != nil {
		return nil, fmt.Errorf("Can't read directory %w", err)
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPage
	}

	rand.Seed(int64(time.Now().Second()))
	n := rand.Intn(len(files))

	file := files[n]

	filename := filepath.Join(file_path, file.Name())

	raw_bits, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Can't open this file %w", err)
	}

	var p storage.Page

	if err := gob.NewDecoder(raw_bits).Decode(&p); err != nil {
		return nil, fmt.Errorf("Can't open decode file %w", err)
	}
	raw_bits.Close()

	return &p, nil
}

func (s Storage) Remove(p *storage.Page) error {
	file_name, err := FileName(p)
	if err != nil {
		return fmt.Errorf("Can't create filename %w", err)
	}
	path := filepath.Join(s.base_path, p.UserName, file_name)
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't delete file %s %w", file_name, err)
		return fmt.Errorf(msg)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	file_name, err := FileName(p)
	if err != nil {
		return false, fmt.Errorf("Can't create filename %w", err)
	}

	path := filepath.Join(s.base_path, p.UserName, file_name)

	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Can't check file: %w", err)
	}
	return true, nil
}

func FileName(p *storage.Page) (string, error) {
	return p.Hash()
}
