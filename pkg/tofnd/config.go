package tofnd

import (
	"fmt"

	"github.com/scalarorg/scalar-healer/config"
)

type ClientConfig struct {
	Address string `json:"address"`
	PartyID string `json:"party_id"`
	KeyUID  string `json:"key_uid"`
	Weight  int    `json:"weight"`
}

func ReadTofndClientConfig(configPath string) ([]ClientConfig, error) {
	cfgPath := fmt.Sprintf("%s/tofnd.json", configPath)
	configs, err := config.ReadJsonArrayConfig[ClientConfig](cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tofnd client configs: %w", err)
	}

	return configs, nil
}
