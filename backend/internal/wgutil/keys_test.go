package wgutil

import (
    "testing"
)

func TestGenerateKeys(t *testing.T) {
    privateKey, publicKey, err := GenerateKeys()
    if err != nil {
        t.Fatalf("GenerateKeys failed: %v", err)
    }

    if len(privateKey) == 0 || len(publicKey) == 0 {
        t.Fatal("Expected non-empty private and public keys")
    }

}
