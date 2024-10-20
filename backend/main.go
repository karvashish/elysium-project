package main

import (
	"elysium-backend/internal/wgutil"
	"fmt"
)

func main() {
	privateKey, publicKey, err := wgutil.GenerateKeys()
	if err != nil {
		fmt.Printf("Failed to generate keys: %v\n", err)
		return
	}

	path := "config/keys"
	filename := "server_private.key"
	if err := wgutil.SaveKeyToFile(path, filename, privateKey); err != nil {
		fmt.Printf("Failed to save private key: %v\n", err)
		return
	}

	fmt.Printf("Keys generated successfully.\nPrivate Key saved as '%s/%s'.\nPublic Key: %s\n", path, filename, publicKey)
}
