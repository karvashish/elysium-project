package services

import (
	"bufio"
	"elysium-backend/config"
	"elysium-backend/internal/models"
	"elysium-backend/internal/repositories"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func InsertPeer(newPeer *models.Peer) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.InsertPeer -> called")
	}

	if newPeer.AssignedIP == nil {
		log.Println("services.InsertPeer -> IP not assigned, requesting new IP")
		if err := assignNewIP(newPeer); err != nil {
			log.Println("services.InsertPeer -> Error assigning new IP:", err)
			return err
		}
	}

	if err := repositories.InsertPeer(newPeer); err != nil {
		log.Println("services.InsertPeer -> Error inserting peer:", err)
		return err
	}

	return nil
}

func assignNewIP(newPeer *models.Peer) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.assignNewIP -> called")
	}
	newPeer.AssignedIP = net.ParseIP("10.0.0.2")
	return nil
}

func GetPeer(peerID *uuid.UUID) (*models.Peer, error) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.GetPeer -> called")
	}

	peer, err := repositories.GetPeer(*peerID)
	if err != nil {
		log.Println("services.GetPeer -> Error retrieving peer:", err)
		return nil, err
	}
	return peer, nil
}

func CompileClient(pubKey string, target models.OSArch) (string, error) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.CompileClient -> called")
	}

	if err := target.Validate(); err != nil {
		log.Println("services.CompileClient -> invalid target:", err)
		return "", err
	}

	clientDir := config.GetEnv("CLIENT_DIR", "./client")
	binaryName := config.GetEnv("BINARY_NAME", "elysium-client")
	if target == "x86_64-pc-windows-gnu" {
		binaryName += ".exe"
	}
	outputDir := config.GetEnv("OUTPUT_DIR", "./compiled_binaries")
	compileArgs := strings.Fields(config.GetEnv("COMPILE_ARGS", ""))

	args := append([]string{"build", "--release", "--target", string(target)}, compileArgs...)
	cmd := exec.Command("cargo", args...)
	cmd.Dir = clientDir
	cmd.Env = append(os.Environ(), "ADDR=10.0.0.2", "CIDR="+fmt.Sprint(24), "SERVERPUB="+pubKey, "SERVERENDPOINT=192.168.0.1:51820", "SERVERIP=10.0.0.1")
	if target == models.OSArchx86_64Linux {
		cmd.Env = append(cmd.Env, "RUSTFLAGS=-C linker=x86_64-linux-gnu-gcc")
	}

	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Println("services.CompileClient -> command start failed:", err)
		return "", err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		log.Println("services.CompileClient -> build output:", scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		log.Println("services.CompileClient -> command execution failed:", err)
		return "", err
	}

	destPath := filepath.Join(outputDir, fmt.Sprintf("%d", time.Now().UnixNano()), binaryName)
	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		log.Println("services.CompileClient -> failed to create directory:", err)
		return "", err
	}

	sourcePath := filepath.Join(clientDir, "target", string(target), "release", binaryName)
	if err := os.Rename(sourcePath, destPath); err != nil {
		log.Println("services.CompileClient -> failed to move executable:", err)
		return "", err
	}

	relativePath, _ := filepath.Rel(outputDir, destPath)
	return relativePath, nil
}
