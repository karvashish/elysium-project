package wgutil

import (
	"fmt"
	"os"

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

func InitWireGuardInterface() error {
	err := CreateWireGuardInterface("wg0")
	if err != nil {
		return fmt.Errorf("failed to create WireGuard interface: %w", err)
	}

	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to initialize WireGuard client: %w", err)
	}
	defer client.Close()

	privateKeyPath := "config/keys/server_private.key"
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := wgtypes.ParseKey(string(privateKeyData))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: new(int),
	}

	err = client.ConfigureDevice("wg0", config)
	if err != nil {
		return fmt.Errorf("failed to configure WireGuard interface: %w", err)
	}

	return nil
}
