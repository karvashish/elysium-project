package config

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Ip_Range struct {
	Start net.IP
	End   net.IP
}

var range_max int64 = 256

var logLevel string

var Ip_ranges []Ip_Range

func setLogLevel() {
	logLevel = GetEnv("LOG_LEVEL", "INFO")
}

func GetLogLevel() string {
	return logLevel
}

func GetIpRanges() []Ip_Range {
	return Ip_ranges
}

func LoadEnv(provided_path string) {
	err := godotenv.Load(provided_path)
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
	setLogLevel()
	Ip_ranges, _ = GenerateIPRanges(GetEnv("BACKEND_WG_IP", "10.0.0.1"), GetEnv("WG_NETWORK_MASK", "/24"))
	for i, r := range Ip_ranges {
		fmt.Printf("Range %d: Start = %s, End = %s\n", i+1, r.Start.String(), r.End.String())
	}

}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GenerateIPRanges(serverIP, mask string) ([]Ip_Range, error) {
	_, ipNet, err := net.ParseCIDR(serverIP + mask)
	if err != nil {
		return nil, err
	}

	ones, bits := ipNet.Mask.Size()
	if bits != 32 {
		return nil, fmt.Errorf("only IPv4 CIDRs are supported")
	}

	total_ips := (1 << (bits - ones)) - 1

	var ranges []Ip_Range

	network_number := new(big.Int).SetBytes(ipNet.IP.To4()).Int64()
	broadcast_number := network_number + int64(total_ips)

	for i := network_number + 1; i < broadcast_number; i += range_max {
		new_start_ip := make([]byte, 4)
		binary.BigEndian.PutUint32(new_start_ip, uint32(i))
		if new_start_ip[3] != 0 {
			i -= 1
		}
		new_end_ip := make([]byte, 4)
		binary.BigEndian.PutUint32(new_end_ip, uint32(math.Min(float64(i+range_max-1), float64(broadcast_number-1))))

		new_range := Ip_Range{Start: new_start_ip, End: new_end_ip}
		ranges = append(ranges, new_range)
	}

	return ranges, nil

}
