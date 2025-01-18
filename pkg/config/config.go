package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BrowserBinary     string            `json:"browserBinary"`
	BrowserArgs       string            `json:"browserArgs"`
	NotificationToken string            `json:"notificationToken"`
	OrgTokens         map[string]string `json:"orgTokens"`
}

func FromFile(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = json.Unmarshal(raw, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
