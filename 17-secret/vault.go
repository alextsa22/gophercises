package _7_secret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/alextsa22/gophercises/17-secret/cipher"
)

type vault struct {
	encodingKey string
	filepath    string
	mutex       *sync.Mutex
	keyValues   map[string]string
}

func NewVault(encodingKey, filepath string) *vault {
	return &vault{
		encodingKey: encodingKey,
		filepath:    filepath,
		mutex:       &sync.Mutex{},
	}
}

func (v *vault) load() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()

	r, err := cipher.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}

	return v.readKeyValues(r)
}

func (v *vault) readKeyValues(r io.Reader) error {
	return json.NewDecoder(r).Decode(&v.keyValues)
}

func (v *vault) save() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := cipher.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}

	return v.writeKeyValues(w)
}

func (v *vault) writeKeyValues(w io.Writer) error {
	return json.NewEncoder(w).Encode(v.keyValues)
}

func (v *vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if err := v.load(); err != nil {
		return "", err
	}

	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}

	return value, nil
}

func (v *vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if err := v.load(); err != nil {
		return err
	}

	v.keyValues[key] = value
	return v.save()
}
