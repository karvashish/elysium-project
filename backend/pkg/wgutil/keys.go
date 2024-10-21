package wgutil

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GenerateKeys() (string, string, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", fmt.Errorf("error generating private key: %w", err)
	}

	publicKey := privateKey.PublicKey()

	return privateKey.String(), publicKey.String(), nil
}

func SaveKeyToFile(path, filename, key string) error {
	keyPath := filepath.Join(path, filename)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(key)
	return err
}
