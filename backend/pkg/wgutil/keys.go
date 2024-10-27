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
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey := privateKey.PublicKey()

	return privateKey.String(), publicKey.String(), nil
}

func SaveKeyToFile(path, filename, key string) error {
	keyPath := filepath.Join(path, filename)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", path, err)
	}

	file, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", keyPath, err)
	}
	defer file.Close()

	_, err = file.WriteString(key)
	if err != nil {
		return fmt.Errorf("failed to write key to file %s: %v", keyPath, err)
	}

	return nil
}
