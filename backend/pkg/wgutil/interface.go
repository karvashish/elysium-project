package wgutil

import (
	"fmt"

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
	err := netlink.LinkAdd(link)
	if err != nil {
		return fmt.Errorf("failed to create WireGuard interface %s: %v", ifaceName, err)
	}

	createdLink, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return fmt.Errorf("failed to get link %s: %v", ifaceName, err)
	}

	err = netlink.LinkSetUp(createdLink)
	if err != nil {
		return fmt.Errorf("failed to bring up link %s: %v", ifaceName, err)
	}

	return nil
}

func setIPAddress(ifaceName string, ipAddress string, ipMask string) error {
	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(ipAddress + ipMask)
	if err != nil {
		return err
	}
	return netlink.AddrAdd(link, addr)
}

func InitWireGuardInterface(server_interface string, server_port int, server_IP string, network_mask string) error {
	err := CreateWireGuardInterface(server_interface)
	if err != nil {
		return fmt.Errorf("failed to create WireGuard interface: %w", err)
	}

	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to initialize WireGuard client: %w", err)
	}
	defer client.Close()

	privKey, pubKey, err := GenerateKeys()
	if err != nil {
		return fmt.Errorf("failed to GenerateKeys: %w", err)
	}
	SaveKeyToFile("config/keys/", "server_private.key", privKey)
	fmt.Println(pubKey)

	privateKey, err := wgtypes.ParseKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &server_port,
	}

	err = client.ConfigureDevice(server_interface, config)
	if err != nil {
		return fmt.Errorf("failed to configure WireGuard interface: %w", err)
	}

	if err := setIPAddress(server_interface, server_IP, network_mask); err != nil {
		return fmt.Errorf("failed to Ip for wg0 interface: %w", err)
	}

	return nil
}
