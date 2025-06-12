package tofnd

import (
	"fmt"

	"github.com/scalarorg/scalar-healer/config"
)

type ClientConfig struct {
	Address string `mapstructure:"address"`
	PartyID string `mapstructure:"party_id"`
	KeyID   string `mapstructure:"key_id"`
	Weight  int    `mapstructure:"weight"`
}

func ReadTofndClientConfig(configPath string) ([]ClientConfig, error) {
	cfgPath := fmt.Sprintf("%s/tofnd.json", configPath)
	configs, err := config.ReadJsonArrayConfig[ClientConfig](cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tofnd client configs: %w", err)
	}
	return configs, nil
}
