package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/k1ntoha/LaterBot/lib/e"
	"github.com/k1ntoha/LaterBot/storage"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0775

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.Wrap("Failed to save", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)
	fmt.Printf("fPath: %v\n", fPath)
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.Wrap("Failed to PickRandom", err) }()

	path := filepath.Join(s.basePath, userName)
	//check user folder
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("Failed to Remove", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)
	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("Failed to remove %s", fName), err)
	}
	return nil

}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("Failed to Check", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(os.ErrNotExist, err):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("Failed to check existence %s", fName), err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("Failed to open file", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("Failed to Decode", err)
	}
	return &p, nil
}
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
