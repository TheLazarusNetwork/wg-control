package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/TheLazarusNetwork/wg-control/model"
	"github.com/TheLazarusNetwork/wg-control/util"
)

// Serialize write interface to disk
func Serialize(id string, c interface{}) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if id == "server.json" {
		return util.WriteFile(filepath.Join(os.Getenv("WG_CONF_DIR"), id), b)
	} else {
		return util.WriteFile(filepath.Join(os.Getenv("WG_KEYS_DIR"), id), b)
	}
}

// Deserialize read interface from disk
func Deserialize(id string) (interface{}, error) {
	var path string
	if id == "server.json" {
		path = filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	} else {
		path = filepath.Join(os.Getenv("WG_KEYS_DIR"), id)
	}

	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if id == "server.json" {
		var s *model.Server
		err = json.Unmarshal(data, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	// if not the server, must be client
	var c *model.Client
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
