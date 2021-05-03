package _7_secret

import (
	"errors"
	"github.com/alextsa22/gophercises/17-secret/encrypt"
)

type vault struct {
	encodingKey string
	keyValues   map[string]string
}

func NewVault(encodingKey string) *vault {
	return &vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
	}
}

func (v *vault) Get(key string) (string, error) {
	encryptedValue, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}

	value, err := encrypt.Decrypt(v.encodingKey, encryptedValue)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (v *vault) Set(key, value string) error {
	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}

	v.keyValues[key] = encryptedValue
	return nil
}
