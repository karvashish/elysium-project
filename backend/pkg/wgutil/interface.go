package wgutil

import (
	"elysium-backend/internal/models"
	"elysium-backend/internal/services"
	"fmt"
	"time"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func CreateWireGuardInterface(ifaceName string) error {
	link := &netlink.GenericLink{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifaceName,
		},
		LinkType: "wireguard",
	}
	if err := netlink.LinkAdd(link); err != nil {
		return fmt.Errorf("error creating WireGuard interface %s: %v", ifaceName, err)
	}

	createdLink, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return fmt.Errorf("error retrieving interface %s after creation: %v", ifaceName, err)
	}

	if err := netlink.LinkSetUp(createdLink); err != nil {
		return fmt.Errorf("error setting interface %s up: %v", ifaceName, err)
	}

	return nil
}

func setIPAddress(ifaceName, ipAddress, ipMask string) error {
	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return fmt.Errorf("error retrieving link for %s: %v", ifaceName, err)
	}

	addr, err := netlink.ParseAddr(ipAddress + ipMask)
	if err != nil {
		return fmt.Errorf("error parsing IP address %s%s: %v", ipAddress, ipMask, err)
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		return fmt.Errorf("error adding IP address %s%s to %s: %v", ipAddress, ipMask, ifaceName, err)
	}

	return nil
}

func InitWireGuardInterface(server_interface string, server_port int, server_IP, network_mask string) error {
	if err := CreateWireGuardInterface(server_interface); err != nil {
		return fmt.Errorf("failed to create WireGuard interface: %v", err)
	}

	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("error initializing WireGuard client: %v", err)
	}
	defer client.Close()

	privKey, pubKey, err := GenerateKeys()
	if err != nil {
		return fmt.Errorf("error generating keys: %v", err)
	}
	SaveKeyToFile("config/keys/", "server_private.key", privKey)

	privateKey, err := wgtypes.ParseKey(privKey)
	if err != nil {
		return fmt.Errorf("error parsing private key: %v", err)
	}

	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &server_port,
	}

	if err := client.ConfigureDevice(server_interface, config); err != nil {
		return fmt.Errorf("error configuring WireGuard interface %s: %v", server_interface, err)
	}

	if err := setIPAddress(server_interface, server_IP, network_mask); err != nil {
		return fmt.Errorf("error setting IP address for interface %s: %v", server_interface, err)
	}

	backend_server := models.Peer{
		PublicKey:  pubKey,
		AssignedIP: server_IP,
		Status:     "active",
		IsGateway:  false,
		CreatedOn:  time.Now().UTC(),
	}

	if err := services.InsertPeer(&backend_server); err != nil {
		return fmt.Errorf("error saving backend server in peer table %v", err)
	}

	return nil
}
