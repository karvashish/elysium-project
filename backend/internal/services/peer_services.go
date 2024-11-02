package services

import (
	"elysium-backend/config"
	"elysium-backend/internal/models"
	"elysium-backend/internal/repositories"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func InsertPeer(newPeer *models.Peer) error {

	if err := repositories.InsertPeer(newPeer); err != nil {
		log.Printf("services.InsertPeer -> Error inserting peer : %v", err)
		return err
	}

	return nil
}

func GetPeer(peerID *uuid.UUID) (*models.Peer, error) {
	peer, err := repositories.GetPeer(*peerID)

	if err != nil {
		log.Printf("services.GetPeer -> Error retrieving peer : %v", err)
		return nil, err
	}
	return peer, nil
}

func CompileClient() (string, error) {
	clientDir := config.GetEnv("CLIENT_DIR", "./client")
	binaryName := config.GetEnv("BINARY_NAME", "elysium-client")
	outputDir := config.GetEnv("OUTPUT_DIR", "./compiled_binaries")

	compileArgs := config.GetEnv("COMPILE_ARGS", "")
	args := append([]string{"build", "--release"}, strings.Fields(compileArgs)...)
	cmd := exec.Command("cargo", args...)
	cmd.Dir = clientDir

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("compilation error: %w", err)
	}

	uniqueDir := filepath.Join(outputDir, fmt.Sprintf("%d", time.Now().UnixNano()))
	if err := os.MkdirAll(uniqueDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	sourcePath := filepath.Join(clientDir, "target", "release", binaryName)
	destPath := filepath.Join(uniqueDir, binaryName)

	if err := os.Rename(sourcePath, destPath); err != nil {
		return "", fmt.Errorf("failed to move executable: %w", err)
	}

	relativePath := filepath.Join(uniqueDir[len(outputDir)+1:], binaryName)
	return relativePath, nil
}
