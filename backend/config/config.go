package config

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Ip_Range struct {
	Start net.IP
	End   net.IP
}

var range_max int64 = 255

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
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GenerateIPRanges(serverIP string, mask string) ([]Ip_Range, error) {
	_, ipNet, err := net.ParseCIDR(serverIP + mask)
	if err != nil {
		return nil, err
	}

	var ipRanges []Ip_Range

	ip := ipNet.IP.To4()
	if ip == nil {
		return nil, fmt.Errorf("invalid IPv4 address: %s", ipNet.IP)
	}

	server := new(big.Int).SetBytes(net.ParseIP(serverIP).To4()).Int64()
	network_number := new(big.Int).SetBytes(ip).Int64()

	ones, bits := ipNet.Mask.Size()
	if bits != 32 {
		return nil, fmt.Errorf("only IPv4 CIDRs are supported")
	}

	avail_ips := (1 << (bits - ones)) - 1
	last_ip_int := network_number + int64(avail_ips)

	i := network_number + 1
	for i < last_ip_int {
		new_start_ip := make([]byte, 4)
		binary.BigEndian.PutUint32(new_start_ip, uint32(i))

		var new_range Ip_Range
		if i < server && server < i+range_max {
			bsi := make([]byte, 4)
			binary.BigEndian.PutUint32(bsi, uint32(server-1))
			new_range_1 := Ip_Range{
				Start: new_start_ip,
				End:   bsi,
			}
			ipRanges = append(ipRanges, new_range_1)

			new_start_ip = make([]byte, 4)
			binary.BigEndian.PutUint32(new_start_ip, uint32(server+1))
		}

		if i == network_number+1 {
			i += range_max - 1
		} else {
			i += range_max
		}

		new_last_ip := make([]byte, 4)
		if i > last_ip_int-1 {
			binary.BigEndian.PutUint32(new_last_ip, uint32(last_ip_int-1))
		} else {
			binary.BigEndian.PutUint32(new_last_ip, uint32(i))
			i += 1
		}

		new_range = Ip_Range{
			Start: net.IP(new_start_ip),
			End:   net.IP(new_last_ip),
		}
		ipRanges = append(ipRanges, new_range)
	}

	println(server, network_number, last_ip_int)
	return ipRanges, nil
}
