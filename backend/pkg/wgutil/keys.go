package wgutil

import (
	"elysium-backend/config"
	"log"
	"os"
	"path/filepath"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GenerateKeys() (string, string, error) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("wgutil.GenerateKeys -> called")
	}

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Println("wgutil.GenerateKeys -> failed to generate private key:", err)
		return "", "", err
	}

	publicKey := privateKey.PublicKey()

	log.Println("wgutil.GenerateKeys -> keys generated successfully")
	return privateKey.String(), publicKey.String(), nil
}

func SaveKeyToFile(path, filename, key string) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("wgutil.SaveKeyToFile -> called with path:", path, "filename:", filename)
	}

	keyPath := filepath.Join(path, filename)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Println("wgutil.SaveKeyToFile -> failed to create directory:", err)
		return err
	}

	file, err := os.Create(keyPath)
	if err != nil {
		log.Println("wgutil.SaveKeyToFile -> failed to create file:", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(key)
	if err != nil {
		log.Println("wgutil.SaveKeyToFile -> failed to write key to file:", err)
		return err
	}

	if config.GetLogLevel() == "DEBUG" {
		log.Println("wgutil.SaveKeyToFile -> key saved successfully to:", keyPath)
	}
	return nil
}
