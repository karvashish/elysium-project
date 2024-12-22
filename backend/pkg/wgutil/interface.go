package wgutil

import (
	"elysium-backend/internal/models"
	"elysium-backend/internal/services"
	"log"
	"net"
	"time"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func CreateWireGuardInterface(ifaceName string) error {
	log.Println("wgutil.CreateWireGuardInterface -> called with ifaceName:", ifaceName)

	link := &netlink.GenericLink{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifaceName,
		},
		LinkType: "wireguard",
	}
	if err := netlink.LinkAdd(link); err != nil {
		log.Println("wgutil.CreateWireGuardInterface -> error creating WireGuard interface:", err)
		return err
	}

	createdLink, err := netlink.LinkByName(ifaceName)
	if err != nil {
		log.Println("wgutil.CreateWireGuardInterface -> error retrieving interface after creation:", err)
		return err
	}

	if err := netlink.LinkSetUp(createdLink); err != nil {
		log.Println("wgutil.CreateWireGuardInterface -> error setting interface up:", err)
		return err
	}

	log.Println("wgutil.CreateWireGuardInterface -> successfully created and set up interface:", ifaceName)
	return nil
}

func setIPAddress(ifaceName, ipAddress, ipMask string) error {
	log.Println("wgutil.setIPAddress -> called with ifaceName:", ifaceName, "ipAddress:", ipAddress, "ipMask:", ipMask)

	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		log.Println("wgutil.setIPAddress -> error retrieving link:", err)
		return err
	}

	addr, err := netlink.ParseAddr(ipAddress + ipMask)
	if err != nil {
		log.Println("wgutil.setIPAddress -> error parsing IP address:", err)
		return err
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		log.Println("wgutil.setIPAddress -> error adding IP address:", err)
		return err
	}

	log.Println("wgutil.setIPAddress -> successfully set IP address:", ipAddress+ipMask)
	return nil
}

func InitWireGuardInterface(server_interface string, server_port int, server_IP net.IP, network_mask string) error {
	log.Println("wgutil.InitWireGuardInterface -> called")

	if err := CreateWireGuardInterface(server_interface); err != nil {
		log.Println("wgutil.InitWireGuardInterface -> failed to create WireGuard interface:", err)
		return err
	}

	client, err := wgctrl.New()
	if err != nil {
		log.Println("wgutil.InitWireGuardInterface -> error initializing WireGuard client:", err)
		return err
	}
	defer client.Close()

	privKey, pubKey, err := GenerateKeys()
	if err != nil {
		log.Println("wgutil.InitWireGuardInterface -> error generating keys:", err)
		return err
	}
	SaveKeyToFile("config/keys/", "server_private.key", privKey)

	privateKey, err := wgtypes.ParseKey(privKey)
	if err != nil {
		log.Println("wgutil.InitWireGuardInterface -> error parsing private key:", err)
		return err
	}

	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &server_port,
	}

	if err := client.ConfigureDevice(server_interface, config); err != nil {
		log.Println("wgutil.InitWireGuardInterface -> error configuring WireGuard interface:", err)
		return err
	}

	if err := setIPAddress(server_interface, server_IP.String(), network_mask); err != nil {
		log.Println("wgutil.InitWireGuardInterface -> error setting IP address for interface:", err)
		return err
	}

	backend_server := models.Peer{
		PublicKey:  pubKey,
		AssignedIP: server_IP,
		Status:     "active",
		IsGateway:  false,
		CreatedOn:  time.Now().UTC(),
	}

	if err := services.InsertPeer(&backend_server); err != nil {
		log.Println("wgutil.InitWireGuardInterface -> error saving backend server in peer table:", err)
		return err
	}

	log.Println("wgutil.InitWireGuardInterface -> successfully initialized WireGuard interface:", server_interface)
	return nil
}
