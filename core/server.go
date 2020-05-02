package core

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/TheLazarusNetwork/wg-control/model"
	"github.com/TheLazarusNetwork/wg-control/storage"
	"github.com/TheLazarusNetwork/wg-control/template"
	"github.com/TheLazarusNetwork/wg-control/util"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// ReadServer object, create default one
func ReadServer() (*model.Server, error) {
	if !util.FileExists(filepath.Join(os.Getenv("WG_CONF_DIR"), "server.json")) {
		server := &model.Server{}

		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		server.PrivateKey = key.String()
		server.PublicKey = key.PublicKey().String()

		server.Endpoint = os.Getenv("WG_DOMAIN")
		server.ListenPort = 51820

		server.Address = make([]string, 0)
		server.Address = append(server.Address, "10.0.0.1/24")

		server.Dns = make([]string, 0)
		server.Dns = append(server.Dns, os.Getenv("WG_DNS"))

		server.PersistentKeepalive = 16
		server.Mtu = 0
		server.PreUp = "echo WireGuard PreUp"
		server.PostUp = "iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE"
		server.PreDown = "echo WireGuard PreDown"
		server.PostDown = "iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE"
		server.SaveConfig = false
		server.Created = time.Now().UTC()
		server.Updated = server.Created

		err = storage.Serialize("server.json", server)
		if err != nil {
			return nil, err
		}

		// server.json was missing, dump wg config after creation
		err = UpdateServerConfigWg()
		if err != nil {
			return nil, err
		}
	}

	c, err := storage.Deserialize("server.json")
	if err != nil {
		return nil, err
	}

	return c.(*model.Server), nil
}

// UpdateServer keep private values from existing one
func UpdateServer(server *model.Server) (*model.Server, error) {
	current, err := storage.Deserialize("server.json")
	if err != nil {
		return nil, err
	}

	// check if server is valid
	errs := server.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("server validation error")
		}
		return nil, errors.New("failed to validate server")
	}

	server.PrivateKey = current.(*model.Server).PrivateKey
	server.PublicKey = current.(*model.Server).PublicKey
	//server.PresharedKey = current.(*model.Server).PresharedKey
	server.Updated = time.Now().UTC()

	err = storage.Serialize("server.json", server)
	if err != nil {
		return nil, err
	}

	v, err := storage.Deserialize("server.json")
	if err != nil {
		return nil, err
	}
	server = v.(*model.Server)

	return server, UpdateServerConfigWg()
}

// UpdateServerConfigWg in wg format
func UpdateServerConfigWg() error {
	clients, err := ReadClients()
	if err != nil {
		return err
	}

	server, err := ReadServer()
	if err != nil {
		return err
	}

	_, err = template.DumpServerWg(clients, server)
	if err != nil {
		return err
	}

	return nil
}

// GetAllReservedIps the list of all reserved IPs, client and server
func GetAllReservedIps() ([]string, error) {
	clients, err := ReadClients()
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	reserverIps := make([]string, 0)

	for _, client := range clients {
		for _, cidr := range client.Address {
			ip, err := util.GetIpFromCidr(cidr)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"cidr": cidr,
				}).Error("failed to ip from cidr")
			} else {
				reserverIps = append(reserverIps, ip)
			}
		}
	}

	for _, cidr := range server.Address {
		ip, err := util.GetIpFromCidr(cidr)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"cidr": err,
			}).Error("failed to ip from cidr")
		} else {
			reserverIps = append(reserverIps, ip)
		}
	}

	return reserverIps, nil
}
