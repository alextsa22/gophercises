package _7_secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/alextsa22/gophercises/17-secret/encrypt"
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

func (v *vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()

	var sb strings.Builder
	if _, err = io.Copy(&sb, f); err != nil {
		return err
	}

	decryptedJson, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	r := strings.NewReader(decryptedJson)
	decoder := json.NewDecoder(r)
	if err = decoder.Decode(&v.keyValues); err != nil {
		return err
	}

	return nil
}

func (v *vault) saveKeyValues() error {
	var sb strings.Builder
	encoder := json.NewEncoder(&sb)
	if err := encoder.Encode(v.keyValues); err != nil {
		return err
	}

	encryptedJson, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = fmt.Fprint(f, encryptedJson); err != nil {
		return err
	}

	return nil
}

func (v *vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if err := v.loadKeyValues(); err != nil {
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

	if err := v.loadKeyValues(); err != nil {
		return err
	}

	v.keyValues[key] = value
	return v.saveKeyValues()
}
