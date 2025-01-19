package services

import (
	"bufio"
	"crypto/sha256"
	"elysium-backend/config"
	"elysium-backend/internal/models"
	"elysium-backend/internal/repositories"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Ip_Range struct {
	name  string
	start net.IP
	end   net.IP
}

var ip_ranges []Ip_Range = []Ip_Range{
	{name: "range1", start: net.IPv4(10, 0, 0, 2), end: net.IPv4(10, 0, 0, 255)},
}

func InsertPeer(newPeer *models.Peer) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.InsertPeer -> called")
	}

	if err := repositories.InsertPeer(newPeer); err != nil {
		log.Println("services.InsertPeer -> Error inserting peer:", err)
		return err
	}

	return nil
}

func timeToIp(time *time.Time) net.IP {
	hash := sha256.Sum256([]byte(time.String()))
	hash_int := new(big.Int).SetBytes(hash[:])

	index := new(big.Int).
		Mod(hash_int, big.NewInt(int64(len(ip_ranges)))).
		Int64()

	selectedRange := ip_ranges[index]

	range_start := new(big.Int).SetBytes(selectedRange.start.To4()).Int64()
	range_end := new(big.Int).SetBytes(selectedRange.end.To4()).Int64()
	range_size := range_end - range_start + 1

	ip_offset := new(big.Int).Mod(hash_int, big.NewInt(range_size)).Int64()
	allocated_ip_int := range_start + ip_offset

	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, uint32(allocated_ip_int))

	allocatedIP := net.IP(ipBytes)
	return allocatedIP
}

func AssignNewIP(newPeer *models.Peer) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.assignNewIP -> called")
	}

	allocatedIP := timeToIp(&newPeer.CreatedOn)

	is_avail, _ := repositories.IsIpAvailable(allocatedIP)

	if is_avail {
		log.Println("services.assignNewIP -> IP Allocated")
		newPeer.AssignedIP = allocatedIP
		return nil
	} else {
		retry := 0
		for retry < 100 {
			retry++
			adjustedTime := newPeer.CreatedOn.Add(time.Duration(retry) * time.Millisecond)
			allocatedIP = timeToIp(&adjustedTime)

			is_avail, _ = repositories.IsIpAvailable(allocatedIP)
			if is_avail {
				log.Println("services.assignNewIP -> IP Allocated")
				newPeer.AssignedIP = allocatedIP
				return nil
			}
		}
		return fmt.Errorf("services.assignNewIP -> Error no available IP after 10 retries")
	}
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

func CompileClient(pubKey string, target models.OSArch, assignedIp net.IP) (string, error) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("services.CompileClient -> called")
	}
	serverIp := config.GetEnv("BACKEND_WG_IP", "10.0.0.1")

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
	cmd.Env = append(
		os.Environ(),
		"ADDR="+assignedIp.String(),
		"CIDR="+fmt.Sprint(24),
		"SERVERPUB="+pubKey,
		"SERVERENDPOINT=192.168.0.1:51820",
		"SERVERIP="+serverIp,
	)
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
